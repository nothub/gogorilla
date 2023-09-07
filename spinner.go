package main

import (
	"fmt"
	"os"
	"time"
)

var moon = []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}

func spin() chan any {
	var done = make(chan any)
	go func() {
		for i := 0; ; i++ {
			select {
			case <-done:
				os.Stderr.WriteString("\033[2K\r")
				return
			default:
				os.Stderr.WriteString(fmt.Sprintf("\r%s", moon[(i%len(moon))]))
				time.Sleep(150 * time.Millisecond)
			}
		}
	}()
	return done
}
