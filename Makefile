.PHONY: run build bench run-go build-go

DAYS=$(shell find . -maxdepth 1 -type d -name 'day*' -printf '%P.target\n' | sort)
RUN_DAYS=$(DAYS:.target=.run)
RUN_GO_DAYS=$(DAYS:.target=.run-go)
BUILD_DAYS=$(DAYS:.target=.build)
BUILD_GO_DAYS=$(DAYS:.target=.build-go)

run: $(RUN_DAYS)

run-go: $(RUN_GO_DAYS)

build-go: $(BUILD_GO_DAYS)

build: $(BUILD_DAYS)

bench: build
	./bench.sh

.PHONY: $(RUN_DAYS)

$(RUN_DAYS): %.run:
	$(MAKE) -C $* -s run

.PHONY: $(RUN_GO_DAYS)

$(RUN_GO_DAYS): %.run-go:
	$(MAKE) -C $* -s run-go

.PHONY: $(BUILD_DAYS)

$(BUILD_DAYS): %.build:
	$(MAKE) -C $* -s build

.PHONY: $(BUILD_GO_DAYS)

$(BUILD_GO_DAYS): %.build-go:
	$(MAKE) -C $* -s build-go
