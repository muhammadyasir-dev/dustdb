run:
	@go run cmd/main.go

debug:
	go test -v ./... --race
