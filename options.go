package main

import (
	"fmt"
	"github.com/google/uuid"
	flag "github.com/spf13/pflag"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var (
	uid           string
	shell         string
	server        string
	printJson     bool
	noInteractive bool
)

const usage = `gorilla - a client for Gorilla-CLI LLM servers

Usage:
  gorilla [OPTIONS] [INPUT...]

Options:
  -u, --uid <string>
      User id [default: random uuid].
  -s, --shell <path>
      Shell path for command execution [default: "bash"].
  --server <host>
      LLM host [default: "cli.gorilla-llm.com"].
  --json
      Print commands as json to stdout, implies --no-interactive.
  --no-interactive
      Skip interactive user input.
  -h, --help
      Print help message and exit.

Examples:
  gorilla -u foobar@example.org
      # Set user id and read input interactively.
  gorilla "open the best editor"
      # Read input from arguments and select result interactively.
  gorilla --json check the weather > cmds.json
      # Store results in json format.
  gorilla --shell sh -- http get example.org with curl flag -L
      # Separate options and arguments`

func init() {
	flag.StringVarP(&uid, "uid", "u", "", "")
	flag.StringVarP(&shell, "shell", "s", "bash", "")
	flag.StringVar(&server, "server", "cli.gorilla-llm.com", "")
	flag.BoolVar(&printJson, "json", false, "")
	flag.BoolVar(&noInteractive, "no-interactive", false, "")
	flag.Usage = func() {
		fmt.Println(usage)
		os.Exit(0)
	}
	flag.Parse()

	if strings.TrimSpace(uid) == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			data, err := os.ReadFile(filepath.Join(home, ".gorilla-cli-userid"))
			if err == nil {
				uid = strings.TrimSpace(string(data))
			}
		}
		if strings.TrimSpace(uid) == "" {
			uid = uuid.New().String()
		}
	}

	_, err := url.Parse(fmt.Sprintf("https://%s/", server))
	if err != nil {
		log.Fatal(err)
	}

	if printJson {
		noInteractive = true
	}
}
