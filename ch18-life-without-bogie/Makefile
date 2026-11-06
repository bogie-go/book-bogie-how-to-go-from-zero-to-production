.PHONY: run build fmt vet check

run:
	go run .

build:
	go build -o bin/without-bogie .

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet
