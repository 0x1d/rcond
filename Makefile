SHELL := bash
ARCH ?= amd64
ADDR ?= 0.0.0.0:8080

generate:
	swagger generate server -f api/rcond.yaml -t api/
	go mod tidy

build:
	mkdir -p bin
	env GOOS=linux GOARCH=${ARCH} go build -o bin/rcond-${ARCH} ./cmd/rcond/main.go

run:
	bin/rcond-${ARCH} -config config.yaml

dev:
	RCOND_ADDR=127.0.0.1:8080 \
	RCOND_API_TOKEN=1234567890 \
	go run cmd/rcond/main.go

upload:
	scp bin/rcond-${ARCH} pi@rpi-test:/home/pi/rcond
