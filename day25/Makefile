.PHONY: run build bench go run-go build-go rs run-rs build-rs

BENCHARGS=-w 8
BENCH=hyperfine
BIN=$(notdir $(CURDIR))

GOBIN=./bin/$(BIN)
RSBIN=./target/release/$(BIN)

GOSRC=$(shell find . -type f -name '*.go')
RSSRC=$(shell find . -type f -name '*.rs')

run:
	go run main.go

build: build-go build-rs

bench: build
	$(BENCH) $(BENCHARGS) -L bin $(GOBIN),$(RSBIN) '{bin}'

go: build-go run-go

run-go:
	$(GOBIN)

build-go: $(GOBIN)

$(GOBIN): $(GOSRC)
	go build -o $(GOBIN)

rs: build-rs run-rs

run-rs:
	$(RSBIN)

build-rs: $(RSBIN)

$(RSBIN): $(RSSRC)
	cargo build --release
