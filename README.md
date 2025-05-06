# rcond

A simple daemon and REST API designed to simplify the management of various system components, including:
- Network connections: Utilizing NetworkManager's D-Bus interface to dynamically configure and monitor network connections
- System hostname: Interacting with the hostname1 service to dynamically update the system's hostname
- Authorized SSH keys: Directly managing the user's authorized_keys file to securely add, remove, or modify authorized SSH keys

## Requirements

- Make
- Go 1.19 or later
- NetworkManager
- systemd
- Linux operating system

## Installation

In order to install `rcond` as a systemd service, you need to specify the target architecture and then run the build and install make targets.

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
| POST    | `/network/ap`                       | Create and activate a WiFi access point |
| PUT     | `/network/interface/{interface}`    | Activate a connection                   |
| DELETE  | `/network/interface/{interface}`    | Deactivate a connection                 |
| DELETE  | `/network/connection/{uuid}`        | Remove a connection                     |
| GET     | `/hostname`                         | Get the hostname                        |
| POST    | `/hostname`                         | Set the hostname                        |
| POST    | `/users/{user}/keys`                | Add an authorized SSH key               |
| DELETE  | `/users/{user}/keys/{fingerprint}`  | Remove an authorized SSH key            |

### Response Codes

- 200: Success
- 400: Bad request (invalid JSON payload)
- 405: Method not allowed
- 500: Internal server error

### Request/Response Format
All endpoints use JSON for request and response payloads.

## Examples

### Setup an Access Point

```bash
#!/usr/bin/env bash
set -euo pipefail

# Example script to create and activate a WiFi access point
# Requires:
#   - RCOND_ADDR       (default: http://0.0.0.0:8080)
#   - RCOND_API_TOKEN  (your API token)

API_URL="${RCOND_ADDR:-http://0.0.0.0:8080}"
API_TOKEN="${RCOND_API_TOKEN:-your_api_token}"
INTERFACE="wlan0"
SSID="MyAccessPoint"
PASSWORD="StrongPassword"

echo "Creating access point '$SSID' on interface '$INTERFACE'..."
ap_response=$(curl -sSf -X POST "$API_URL/network/ap" \
  -H "Content-Type: application/json" \
  -H "X-API-Token: $API_TOKEN" \
  -d '{
    "interface": "'"$INTERFACE"'",
    "ssid": "'"$SSID"'",
    "password": "'"$PASSWORD"'"
  }')

# Extract the UUID from the JSON response
AP_UUID=$(echo "$ap_response" | jq -r '.uuid')

echo "Activating connection with UUID '$AP_UUID' on interface '$INTERFACE'..."
curl -sSf -X PUT "$API_URL/network/interface/$INTERFACE" \
  -H "Content-Type: application/json" \
  -H "X-API-Token: $API_TOKEN" \
  -d '{
    "uuid": "'"$AP_UUID"'"
  }'

echo "Access point '$SSID' is now up and running on $INTERFACE."
```