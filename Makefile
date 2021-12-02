DAYS=$(shell find . -maxdepth 1 -type d -name 'day*' -printf '%P.target\n' | sort)
BUILD_DAYS=$(DAYS:.target=.build)

.PHONY: build bench

build: $(BUILD_DAYS)

bench: build
	./bench.sh

.PHONY: $(BUILD_DAYS)

$(BUILD_DAYS): %.build:
	$(MAKE) -C $* -s build
