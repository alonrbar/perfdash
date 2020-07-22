MAIN    := cmd/perfdash/main.go

default: build

build:
	go build -o out/perfdash.exe ${MAIN}

run:
	go run ${MAIN}