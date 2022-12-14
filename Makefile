.PHONY: all build check test clean

all: clean check build test

clean:
	rm -rf build

check:
	cd scripts && chmod +x run_linters.sh && ./run_linters.sh

memcheck: build
	cd scripts && chmod +x run_memcheck.sh && ./run_memcheck.sh

build:
	cd scripts && chmod +x build.sh && ./build.sh

test: build
	cd scripts && chmod +x run_tests.sh && ./run_tests.sh