# rcond

A simple daemon and REST API to manage:
- network connections through NetworkManager's D-Bus interface
- system hostname through the hostname1 service
- authorized SSH keys through the user's authorized_keys file

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
| POST | `/authorized-key` | Add an authorized SSH key |
| DELETE | `/authorized-key` | Remove an authorized SSH key |

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
  -H "X-API-Token: 1234567890" \
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
  -H "X-API-Token: 1234567890" \
  -d '{
    "interface": "wlan0"
  }'
```

### Remove the stored connection

```bash
curl -v -X POST http://localhost:8080/network/remove \
  -H "X-API-Token: 1234567890" \
  -d '{
    "interface": "wlan0"
  }'
```

### Get the hostname

```bash
curl -v http://localhost:8080/hostname \
  -H "X-API-Token: 1234567890"
```

### Set the hostname

```bash
curl -v -X POST http://localhost:8080/hostname \
  -H "Content-Type: application/json" \
  -H "X-API-Token: 1234567890" \
  -d '{
    "hostname": "MyHostname"
  }'
```

### Add an authorized SSH key

```bash
curl -v -X POST http://localhost:8080/authorized-key \
  -H "Content-Type: application/json" \
  -H "X-API-Token: 1234567890" \
  -d '{
    "user": "pi",
    "pubkey": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC1234567890"
  }'
```

### Remove an authorized SSH key

```bash
curl -v -X DELETE http://localhost:8080/authorized-key \
  -H "Content-Type: application/json" \
  -H "X-API-Token: 1234567890" \
  -d '{
    "user": "pi",
    "pubkey": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC1234567890"
    }'
```