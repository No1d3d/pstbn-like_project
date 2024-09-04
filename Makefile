.PHONY: run-server run-server-release run-shell run-release 

run-server:
	@go run ./cmd/server/main.go

run-server-release:
	@GIN_MODE=release go run ./cmd/server/main.go
		
run-shell:
	@go run ./cmd/cli/main.go

run: run-server

run-release: run-server-release
