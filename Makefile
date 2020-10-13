.PHTON: all build fmt run

all: build fmt run

Binary="SkipTable"

build:
	@go build -o ${Binary} .

fmt:
	@go fmt ./

run:
	@./${Binary}
