SHELL := bash
ARCH ?= amd64
ADDR ?= 0.0.0.0:8080

default: build

.PHONY: info
info:
	@echo "Available targets:"
	@grep '^##' Makefile | sed 's/^##//'

## generate: Generate server from Swagger specs
generate:
	swagger generate server -f api/rcond.yaml -t api/
	go mod tidy

## test: run tests
test:
	go test -v ./...

## build: build binary for target $ARCH
build:
	mkdir -p bin
	env GOOS=linux GOARCH=${ARCH} go build -o bin/rcond-${ARCH} ./cmd/rcond/main.go

## install: build and install binary for target $ARCH as systemd service
install: build
	sudo mkdir -p /etc/rcond
	sudo mkdir -p /var/rcond
	sudo cp config/rcond.yaml /etc/rcond/config.yaml
	sudo cp bin/rcond-${ARCH} /usr/local/bin/rcond
	sudo cp systemd/rcond.service /etc/systemd/system/rcond.service
	sudo systemctl daemon-reload
	sudo systemctl enable rcond
	sudo systemctl start rcond

## uninstall: uninstall systemd service
uninstall:
	sudo systemctl stop rcond
	sudo systemctl disable rcond
	sudo rm -rf /etc/rcond
	sudo rm -rf /var/rcond
	sudo rm /usr/local/bin/rcond
	sudo rm /etc/systemd/system/rcond.service

## run: run and build binary for target $ARCH
run: build
	bin/rcond-${ARCH} -config config/rcond.yaml

## dev: run go programm
dev:
	go run cmd/rcond/main.go -config config/rcond.yaml

## dev-agent: run go programm with agent config
dev-agent:
	go run cmd/rcond/main.go -config config/rcond-agent.yaml

## upload: upload binary of given $ARCH to rpi-test
upload:
	scp bin/rcond-${ARCH} pi@rpi-test:/home/pi/rcond
