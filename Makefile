.PHONY: run build bench

DAYS=$(shell find . -maxdepth 1 -type d -name 'day*' -printf '%P.target\n' | sort)
RUN_DAYS=$(DAYS:.target=.run)
BUILD_DAYS=$(DAYS:.target=.build)

run: $(RUN_DAYS)

build: $(BUILD_DAYS)

bench: build
	./bench.sh

.PHONY: $(RUN_DAYS)

$(RUN_DAYS): %.run:
	$(MAKE) -C $* -s run

.PHONY: $(BUILD_DAYS)

$(BUILD_DAYS): %.build:
	$(MAKE) -C $* -s build
