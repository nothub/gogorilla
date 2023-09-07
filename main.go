package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/samber/lo"
	flag "github.com/spf13/pflag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var httpClient = http.DefaultClient

type InputData struct {
	IID   string `json:"interaction_id"`
	UID   string `json:"user_id"`
	Input string `json:"user_input"`
}

type ResultData struct {
	IID  string `json:"interaction_id"`
	UID  string `json:"user_id"`
	Cmd  string `json:"command"`
	Exit string `json:"exit_condition"`
}

func main() {
	log.SetFlags(0)
	run()
}

func run() {
	input, err := readInput(flag.Args())
	if err != nil {
		log.Fatalln(err.Error())
	}
	if input == "" {
		log.Fatalln(errors.New("missing input"))
	}

	var iid = uuid.New().String()
	cmds, err := fetchCmds(iid, input)
	if err != nil {
		log.Fatalln(err.Error())
	}

	cmds = lo.Filter(cmds, func(item string, index int) bool {
		if item == (": #Do nothing") {
			return false
		}

		if s, _, _ := strings.Cut(item, "#"); strings.TrimSpace(s) == "" {
			return false
		}

		return true
	})

	if printJson {
		if len(cmds) < 1 {
			log.Println("no commands found")
			os.Exit(1)
		}
		var buf bytes.Buffer
		err := json.NewEncoder(&buf).Encode(struct {
			Cmds []string `json:"cmds"`
		}{cmds})
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Print(buf.String())
	} else {
		os.Stderr.WriteString("\n")
		for i, s := range cmds {
			log.Printf("# [%v]\n%s\n\n", i, s)
		}
	}

	if noInteractive {
		return
	}

	os.Stderr.WriteString("exec cmd [0]: ")
	r := bufio.NewReader(os.Stdin)
	s, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	s = strings.TrimSuffix(s, "\n")
	if strings.TrimSpace(s) == "" {
		s = "0"
	}
	idx, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	if idx < 0 || idx >= len(cmds) {
		log.Fatalln("invalid selection")
	}

	shCmd := exec.Command(shell, "-c", cmds[idx])
	shCmd.Stdin = os.Stdin
	shCmd.Stdout = os.Stdout
	shCmd.Stderr = os.Stderr

	os.Stderr.WriteString("\n---- command output start ----\n")
	var exitCode int
	if err := shCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.Success() {
				exitCode = 0
			} else {
				exitCode = exitErr.ExitCode()
			}
		} else {
			log.Fatalln(err.Error())
		}
	}
	os.Stderr.WriteString("\n----- command output end -----\n")

	err = reportResult(iid, cmds[idx], exitCode)
	if err != nil {
		log.Printf("unable to report result: %s\n", err.Error())
	}
}

func readInput(args []string) (string, error) {
	input := strings.Join(args, " ")
	input = strings.TrimSpace(input)
	if input != "" {
		log.Printf("🦍 %s\n", input)
		return input, nil
	}

	if noInteractive {
		return "", nil
	}

	os.Stderr.WriteString("🦍 ")
	r := bufio.NewReader(os.Stdin)
	var err error
	input, err = r.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)

	return input, nil
}

func fetchCmds(interactionId string, input string) ([]string, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(InputData{
		IID:   interactionId,
		UID:   uid,
		Input: input,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/commands", server), &buf)
	if err != nil {
		return nil, err
	}

	res, err := send(req)
	if err != nil {
		return nil, err
	}

	var cmds []string
	err = json.NewDecoder(res.Body).Decode(&cmds)
	if err != nil {
		return nil, err
	}

	return cmds, err
}

func reportResult(interactionId string, cmd string, exitCode int) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(ResultData{
		UID:  uid,
		Cmd:  cmd,
		Exit: strconv.Itoa(exitCode),
		IID:  interactionId,
	})
	if err != nil {
		return err
	}

	log.Printf("\nreporting exit code %v back to server\n", exitCode)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/command-execution-result", server), &buf)
	if err != nil {
		return err
	}

	_, err = send(req)
	if err != nil {
		return err
	}

	return nil
}

func send(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "gogorilla")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Close = true

	done := spin()
	res, err := httpClient.Do(req)
	done <- true
	if err != nil {
		return nil, err
	}

	return res, nil
}
