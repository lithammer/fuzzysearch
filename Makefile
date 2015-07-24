default: build

build:
	go build

test:
	go test -v ./...

benchmark:
	go test -v ./... -bench=.
