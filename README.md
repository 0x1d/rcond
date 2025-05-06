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

### Environment Variables

| Environment Variable | Description                             | Default       |
|----------------------|-----------------------------------------|---------------|
| RCOND_ADDR           | Address to bind the HTTP server to.     | 0.0.0.0:8080  |
| RCOND_API_TOKEN      | API token to use for authentication.    | N/A           |
