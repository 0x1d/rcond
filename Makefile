SHELL := bash
ARCH ?= amd64
ADDR ?= 0.0.0.0:8080

default: build

generate:
	swagger generate server -f api/rcond.yaml -t api/
	go mod tidy

build:
	mkdir -p bin
	env GOOS=linux GOARCH=${ARCH} go build -o bin/rcond-${ARCH} ./cmd/rcond/main.go

install: build
	sudo mkdir -p /etc/rcond
	sudo mkdir -p /var/rcond
	sudo cp config/rcond.yaml /etc/rcond/config.yaml
	sudo cp bin/rcond-${ARCH} /usr/local/bin/rcond
	sudo cp systemd/rcond.service /etc/systemd/system/rcond.service
	sudo systemctl daemon-reload
	sudo systemctl enable rcond
	sudo systemctl start rcond

uninstall:
	sudo systemctl stop rcond
	sudo systemctl disable rcond
	sudo rm -rf /etc/rcond
	sudo rm -rf /var/rcond
	sudo rm /usr/local/bin/rcond
	sudo rm /etc/systemd/system/rcond.service

run: build
	bin/rcond-${ARCH} -config config/rcond.yaml

dev:
	RCOND_ADDR=127.0.0.1:8080 \
	RCOND_API_TOKEN=1234567890 \
	go run cmd/rcond/main.go

upload:
	scp bin/rcond-${ARCH} pi@192.168.1.43:/home/pi/rcond
