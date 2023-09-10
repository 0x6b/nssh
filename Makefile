GO=go
BIN=nssh
SRC=$(shell find . -type f -name '*.go')

$(BIN): $(SRC)
	$(GO) build -trimpath ./cmd/nssh

snapshot:
	which goreleaser && goreleaser --snapshot --skip-publish --clean

clean:
	rm -fr $(BIN) dist

.PHONY: clean
