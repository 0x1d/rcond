ARCH ?= arm64
ADDR ?= 0.0.0.0:8080

generate:
	swagger generate server -f api/rcond.yaml -t api/
	go mod tidy

build:
	mkdir -p bin
	env GOOS=linux GOARCH=${ARCH} go build -o bin/rcond-${ARCH} ./cmd/rcond/main.go

run:
	bin/rcond-${ARCH} ${ADDR}

dev:
	go run cmd/rcond/main.go ${ADDR}

upload:
	scp rcond-${ARCH} pi@rpi-40ac:/home/pi/rcond
