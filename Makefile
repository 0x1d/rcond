SHELL := bash
ARCH ?= arm64
ADDR ?= 0.0.0.0:8080

generate:
	swagger generate server -f api/rcond.yaml -t api/
	go mod tidy

build:
	mkdir -p bin
	env GOOS=linux GOARCH=${ARCH} go build -o bin/rcond-${ARCH} ./cmd/rcond/main.go

run:
	source .env && bin/rcond-${ARCH} ${ADDR}

dev:
	RCOND_API_TOKEN=1234567890 go run cmd/rcond/main.go

upload:
	scp rcond-${ARCH} pi@rpi-40ac:/home/pi/rcond
