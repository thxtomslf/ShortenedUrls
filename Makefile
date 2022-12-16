.PHONY: all build test clean

all: clean build test

clean:
	rm main

build:
	cd scripts && chmod +x build.sh && sh build.sh && ./main

test: build
	cd scripts && chmod +x run_tests.sh && sh run_tests.sh