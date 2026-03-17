BIN     := ratatosk
MODULE  := github.com/kvitrvn/ratatosk
CMD     := ./cmd/ratatosk

LDFLAGS := -s -w

.PHONY: build run test vet clean cross

build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) $(CMD)

run:
	go run $(CMD)

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f $(BIN) dist/

cross:
	GOOS=linux   GOARCH=amd64  go build -ldflags "$(LDFLAGS)" -o dist/$(BIN)-linux-amd64   $(CMD)
	GOOS=linux   GOARCH=arm64  go build -ldflags "$(LDFLAGS)" -o dist/$(BIN)-linux-arm64   $(CMD)
	GOOS=darwin  GOARCH=amd64  go build -ldflags "$(LDFLAGS)" -o dist/$(BIN)-darwin-amd64  $(CMD)
	GOOS=darwin  GOARCH=arm64  go build -ldflags "$(LDFLAGS)" -o dist/$(BIN)-darwin-arm64  $(CMD)
	GOOS=windows GOARCH=amd64  go build -ldflags "$(LDFLAGS)" -o dist/$(BIN)-windows-amd64.exe $(CMD)
