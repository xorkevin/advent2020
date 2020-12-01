.PHONY: run bench go run-go build-go rs run-rs build-rs

BENCHARGS=-w 3
BENCH=hyperfine
BIN=$(notdir $(CURDIR))

GOBIN=./bin/$(BIN)
RSBIN=./target/release/$(BIN)

GOSRC=$(shell find . -type f -name '*.go')
RSSRC=$(shell find . -type f -name '*.rs')

run:
	go run main.go

bench: build-go build-rs
	$(BENCH) $(BENCHARGS) $(GOBIN)
	$(BENCH) $(BENCHARGS) $(RSBIN)

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
