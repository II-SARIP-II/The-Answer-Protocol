BIN_DIR    = bin
SERVER_SRC = ./cmd/server
CLI_SRC    = ./cmd/client_cli
GUI_SRC    = ./cmd/client_gui

SERVER_BIN = $(BIN_DIR)/server
CLI_BIN    = $(BIN_DIR)/client_cli
GUI_BIN    = $(BIN_DIR)/client_gui

all: build

install:
	@go mod tidy

build:
	@go build -o $(BIN_DIR)/ ./cmd/...

server:
	@go build -o $(SERVER_BIN) $(SERVER_SRC)

client_cli:
	@go build -o $(CLI_BIN) $(CLI_SRC)

client_gui:
	@go build -o $(GUI_BIN) $(GUI_SRC)

run_server:
	@go run $(SERVER_SRC)

run_cli:
	@go run $(CLI_SRC)

run_gui:
	@go run $(GUI_SRC)

lint:
	@echo "Checking with go vet..."
	@go vet ./...

lint-strict:
	@echo "Checking formatting with gofmt..."
	@test -z "$$(gofmt -l .)" || (echo "Unformatted files found (run 'make fmt' to fix):" && gofmt -l . && exit 1)
	@echo "Running go vet..."
	@go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		GOTOOLCHAIN=go1.23.6 golangci-lint run ./...; \
	else \
		echo "golangci-lint is not installed (run 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest' to install it)"; \
	fi

fmt:
	@echo "Automatic project reformatting..."
	@gofmt -w .

clean:
	@rm -rf $(BIN_DIR)

fclean: clean

re: fclean all

.PHONY: all install build server client_cli client_gui run_server run_cli run_gui lint lint-strict fmt clean fclean re

