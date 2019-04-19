all: run

run:
	go run .

.PHONY: run

test:
	go test ./...

.PHONY: test
