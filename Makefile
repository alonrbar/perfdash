MAIN    := cmd/perfdash/main.go

default: build

build:
	go build ${MAIN}

run:
	go run ${MAIN}