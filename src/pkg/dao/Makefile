.PHONY : build test clean
all: build test clean

gotest:
	@go test -v -cover

build:
	@go build .

clean:
	@rm -f test.sqlite

test: gotest clean
