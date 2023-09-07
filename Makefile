BIN_NAME := gorilla
BIN_PATH := $(CURDIR)/out/$(BIN_NAME)

GOFLAGS  ?= -tags netgo,timetzdata
LDFLAGS  ?= -ldflags="-extldflags=-static"

$(BIN_PATH): go.mod go.sum $(wildcard *.go)
	go build $(GOFLAGS) $(LDFLAGS) -o $(BIN_PATH)

.PHONY: release
release: clean check
	GOOS=linux  GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BIN_PATH)-linux
	GOOS=linux  GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(BIN_PATH)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(BIN_PATH)-darwin
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(BIN_PATH)-darwin-arm64

.PHONY: clean
clean:
	go clean
	-rm -rf $(shell dirname $(BIN_PATH))

.PHONY: check
check:
	go vet
	go test

README.md: $(BIN_PATH)
	@echo "# gogorilla 🦍" > README.md
	@echo "" >> README.md
	@echo "A client for [gorilla-cli](https://github.com/gorilla-llm/gorilla-cli)," >> README.md
	@echo "the [Gorilla LLM](https://gorilla.cs.berkeley.edu/) shell command generator." >> README.md
	@echo "" >> README.md
	@echo "---" >> README.md
	@echo "" >> README.md
	@echo "\`\`\`" >> README.md
	$(BIN_PATH) --help | tail -n +3  >> README.md
	@echo "\`\`\`" >> README.md
