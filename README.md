# gogorilla 🦍

A client for [gorilla-cli](https://github.com/gorilla-llm/gorilla-cli),
the [Gorilla LLM](https://gorilla.cs.berkeley.edu/) shell command generator.

---

```
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
      # Separate options and arguments
```
