.PHONY: run-server run-shell run

run-server:
	@go run ./cmd/server/main.go

run-shell:
	@go run ./cmd/cli/main.go

run: run-server
