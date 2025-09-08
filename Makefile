BINARY_NAME=commit-hook
HOOK_NAME=prepare-commit-msg
HOOK_SOURCE_DIR=scripts
HOOK_DEST_DIR_GLOBAL=$(HOME)/.git-hooks-global
BINARY_DEST_DIR_GLOBAL?=$(HOME)/.local/bin
BASE_CONFIG_FILE=data/conventional_commits.toml
BASE_CONFIG_DIR=$(HOME)/.commit_hooks

all: help

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	go clean
	rm -f $(BINARY_NAME)

install: build
	@echo "Installing git hook and binary globally..."
	@echo "Installing binary to $(BINARY_DEST_DIR_GLOBAL)..."
	install -d "$(BINARY_DEST_DIR_GLOBAL)"
	install -m 755 $(BINARY_NAME) "$(BINARY_DEST_DIR_GLOBAL)/"
	@echo "Installing hook script..."
	mkdir -p "$(HOOK_DEST_DIR_GLOBAL)"
	cp "$(HOOK_SOURCE_DIR)/$(HOOK_NAME)" "$(HOOK_DEST_DIR_GLOBAL)/$(HOOK_NAME)"
	chmod +x "$(HOOK_DEST_DIR_GLOBAL)/$(HOOK_NAME)"
	git config --global core.hooksPath "$(HOOK_DEST_DIR_GLOBAL)"
	@echo "Installing base conventional_commits.toml"
	cp $(BASE_CONFIG_FILE) $(BASE_CONFIG_DIR)
	@echo "Hook and binary installed globally."
	@echo "Please ensure '$(BINARY_DEST_DIR_GLOBAL)' is in your shell's PATH."

test:
	go test -v ./...

.PHONY: all build run clean install test