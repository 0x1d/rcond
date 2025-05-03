# rcond

A simple daemon to manage network connections through NetworkManager's D-Bus interface.

It provides a REST API to:
- Create and activate WiFi connections
- Deactivate WiFi connections 
- Remove stored connection profiles

The daemon is designed to run on Linux systems with NetworkManager.

## Build and Run

```bash
make build
make run
```

## API

The full API specification can be found in [api/rcond.yaml](api/rcond.yaml).

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check endpoint that returns status |
| POST | `/network/up` | Create and activate a WiFi access point |
| POST | `/network/down` | Deactivate a WiFi interface |
| POST | `/network/remove` | Remove the stored connection profile |
| GET | `/hostname` | Get the hostname |
| POST | `/hostname` | Set the hostname |

### Response Codes

- 200: Success
- 400: Bad request (invalid JSON payload)
- 405: Method not allowed
- 500: Internal server error

### Request/Response Format
All endpoints use JSON for request and response payloads.

### Bring a network up

```bash
curl -v -X POST http://localhost:8080/network/up \
  -H "Content-Type: application/json" \
  -d '{
    "interface": "wlan0",
    "ssid":      "MyNetworkSSID",
    "password":  "SuperSecretPassword"
  }'
```

### Bring a network down

```bash
curl -v -X POST http://localhost:8080/network/down \
  -H "Content-Type: application/json" \
  -d '{
    "interface": "wlan0"
  }'
```

### Remove the stored connection

```bash
curl -v -X POST http://localhost:8080/network/remove
```

### Get the hostname

```bash
curl -v http://localhost:8080/hostname
```

### Set the hostname

```bash
curl -v -X POST http://localhost:8080/hostname \
  -H "Content-Type: application/json" \
  -d '{
    "hostname": "MyHostname"
  }'
```