# rcond

[![main](https://github.com/0x1d/rcond/actions/workflows/main.yaml/badge.svg)](https://github.com/0x1d/rcond/actions/workflows/main.yaml)

A simple daemon and REST API designed to simplify the management of various system components, including:
- Network connections: Utilizing NetworkManager's D-Bus interface to dynamically configure network connections
- System hostname: Dynamically update the system's hostname
- Authorized SSH keys: Directly managing the user's authorized_keys file to securely add, remove, or modify authorized SSH keys
- System state: Restart and shutdown the system

## Requirements

- Make
- Go
- NetworkManager
- systemd
- Linux operating system

## Installation

In order to install `rcond` as a systemd service, you need to specify the target architecture and then run the install make target.

```sh
export ARCH=arm64
make install
```

## Run

The run target will build the binary for target architecture and runs it using the default configuration in `config/rcond.yaml`

```sh
make run
```

## Develop

The dev target will run the main.go directly with environment variable configuration:
- RCOND_ADDR = 127.0.0.1:8080
- RCOND_API_TOKEN = 1234567890

```sh
make dev
```

## Configuration

### File

The default config file location is `/etc/rcond/config.yaml`.  
It can be overwritten by environment variables and flags.  
An full example configuration with comments can be found in `config/rcond.yaml`

Example configuration:
```yaml
rcond:
  addr: 0.0.0.0:8080
  api_token: 1234567890
```

### Environment Variables

| Environment Variable | Description                             | Default       |
|----------------------|-----------------------------------------|---------------|
| RCOND_ADDR           | Address to bind the HTTP server to.     | 0.0.0.0:8080  |
| RCOND_API_TOKEN      | API token to use for authentication.    | N/A           |

## API

The full API specification can be found in [api/rcond.yaml](api/rcond.yaml).

### Authentication

All endpoints except `/health` require authentication via an API token passed in the `X-API-Token` header. The token is configured via the `RCOND_API_TOKEN` environment variable when starting the daemon.

### Endpoints
| Method  | Path                                | Description                             |
|---------|-------------------------------------|-----------------------------------------|
| GET     | `/health`                           | Health check endpoint                   |
| POST    | `/network/ap`                       | Create a WiFi access point              |
| POST    | `/network/sta`                      | Connect to a WiFi access point          |
| PUT     | `/network/interface/{interface}`    | Activate a connection                   |
| DELETE  | `/network/interface/{interface}`    | Deactivate a connection                 |
| DELETE  | `/network/connection/{uuid}`        | Remove a connection                     |
| GET     | `/hostname`                         | Get the hostname                        |
| POST    | `/hostname`                         | Set the hostname                        |
| POST    | `/users/{user}/keys`                | Add an authorized SSH key               |
| DELETE  | `/users/{user}/keys/{fingerprint}`  | Remove an authorized SSH key            |
| POST    | `/system/restart`                   | Restart the system                      |
| POST    | `/system/shutdown`                  | Shutdown the system                     |

### Response Codes

- 200: Success
- 400: Bad request (invalid JSON payload)
- 405: Method not allowed
- 500: Internal server error

### Request/Response Format
All endpoints use JSON for request and response payloads.

## Examples

### Connect to a WiFi Access Point

This example will automatically connect to a WiFi access point with the given SSID and password on the interface "wlan0".

```bash
curl -sSf -X POST "http://rpi-test:8080/network/sta" \
  -H "Content-Type: application/json" \
  -H "X-API-Token: 1234567890" \
  -d '{
    "interface": "wlan0",
    "ssid": "MyAccessPoint",
    "password": "StrongPassword",
    "autoconnect": true
  }'
```

### Setup an Access Point

This example will create an access point on the interface "wlan0" with the given SSID and password.

```bash
curl -sSf -X POST "http://rpi-test:8080/network/ap" \
  -H "Content-Type: application/json" \
  -H "X-API-Token: 1234567890" \
  -d '{
    "interface": "wlan0",
    "ssid": "MyAccessPoint",
    "password": "StrongPassword",
    "autoconnect": true
  }'
```
