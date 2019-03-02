export GO111MODULE = on

.PHONY: default test test-cover dev

# for dev
dev:
	fresh

build:
	packr2
	go build -tags netgo -o forest 

clean:
	packr2 clean

release:
	go mod tidy