.PHONY: help run-tui run-cli run-server default golint

default: help

help:
	@echo "usage: make [target]"
	@echo "targets:"
	@echo "	run-server		Start server and begin scraping events"
	@echo "	run-cli			Start command line user interface"
	@echo "	run-tui			Start terminal user interface"
	@echo "	golint			Run golangci-lint check"
	@echo "	precomm			Run precommit check"

run-tui:
	go run ./cmd/tui/

run-cli:
	go run ./cmd/cli/

run-server: 
	go run ./cmd/server/

golint: 
	golangci-lint run

precomm:
	pre-commit run --all-files