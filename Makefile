testAppExe := test-warden.exe

all: run

run:
	-rm -f $(testAppExe)
	go build -o $(testAppExe) .
	$(testAppExe)
# Removes application binary if possible. Because we usually use CTRL-C to kill the process
# and because there is no way for make to trap this signal, we usually are not able to
# remove the binary (running a HTTP server) on exit.
	-rm -f $(testAppExe)

.PHONY: run

test:
	go test ./...

.PHONY: test
