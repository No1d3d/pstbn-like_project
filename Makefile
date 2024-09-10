.PHONY: run-server run-server-release run-release 

OUTPUT := ./build/srv
MAIN_FILE := ./cmd/server/main.go

build:
	@go build -o $(OUTPUT) $(MAIN_FILE)

build-release:
	@GIN_MODE=release go build -o $(OUTPUT) $(MAIN_FILE)

run: build
	./$(OUTPUT)

run-release: build-release
	./$(OUTPUT)
