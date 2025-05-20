# rcond

[![main](https://github.com/0x1d/rcond/actions/workflows/main.yaml/badge.svg)](https://github.com/0x1d/rcond/actions/workflows/main.yaml)

A distributed management daemon designed to remotely configure system components, including:
- Network connections: Manage network connections through NetworkManager's D-Bus interface
- Files: Manage files on the system
- System hostname: Update the system's hostname
- Authorized SSH keys: Manage the user's authorized_keys file to add, remove, or modify authorized SSH keys
- System state: Restart and shutdown the system
- Cluster: Join and manage a cluster of rcond nodes

## Requirements

- Make
- Go
- NetworkManager
- systemd
- Linux

## Installation

In order to build and install `rcond` as a systemd service, you need to specify the target architecture and then run the install make target.

```sh
export ARCH=arm64
make install
```

## Run

To run `rcond` manually, execute the following command in your terminal:
```sh
rcond -config config/rcond.yaml
```

## Development

There are several make targets available:

```text
Available targets:
 generate: Generate server from Swagger specs
 test: run tests
 build: build binary for target $ARCH
 install: build and install binary for target $ARCH as systemd service
 uninstall: uninstall systemd service
 run: run and build binary for target $ARCH
 dev: run go programm
 dev-agent: run go programm with agent config
 upload: upload binary of given $ARCH to rpi-test
```

## Configuration

The default config file location is `/etc/rcond/config.yaml`.  
It can be overwritten by environment variables and flags.  
An full example configuration with comments can be found in `config/rcond.yaml`

### API Server

The API server is the main component of the rcond daemon. It is responsible for managing the host and providing a REST API for managing the system.

Example configuration:
```yaml
rcond:
  addr: 0.0.0.0:8080
  api_token: 1234567890
```

### Network

Network connections can be configured in the `rcond.yaml` file, and these configurations are applied automatically when the node starts up. This allows for easy management of network settings, including the creation of access points and the sharing of network connections, without requiring manual intervention after each reboot.

Here is an example for creating an access point and share network connection on wlan0:

```yaml
network:
  connections:
    # create access point and share network connection on wlan0
    - name: MyHomeWiFi
      id: MyHomeWiFi
      uuid: 222b4580-3e08-4a2c-ae5e-316bb45d44f0
      type: 802-11-wireless
      interface: wlan0
      ssid: MyHomeWiFi
      mode: ap
      band: bg
      channel: 1
      keymgmt: wpa-psk
      psk: SuperSecure
      ipv4method: shared
      ipv6method: ignore
      autoconnect: true
```

### Cluster

The cluster agent is a component of rcond that is responsible for joining and managing a cluster of rcond nodes.
This functionality can be used to manage and configure multiple hosts through a single API server.  
In the background, the cluster agent will use [Serf](https://github.com/hashicorp/serf) to form a cluster, handle node discovery and gossip.

Forming a cluster is optional and can be enabled by configuring the cluster section in the config file.

Example configuration:
```yaml
cluster:
  # Enable the cluster agent 
  enabled: true
  # Name of the node in the cluster
  node_name: rcond
  # Secret key for the cluster agent used for message encryption.
  # Must be 32 bytes long and base64 encoded.
  # Generate with: base64 /dev/urandom | tr -d '\n' | head -c 32
  secret_key: DMXnaJUUbIBMj1Df0dPsQY+Sks1VxWTa
  # Advertise address for the cluster agent
  advertise_addr: 0.0.0.0
  # Advertise port for the cluster agent
  advertise_port: 7946
  # Bind address for the cluster agent
  bind_addr: 0.0.0.0
  # Bind port for the cluster agent
  bind_port: 7946
  # Join addresses for the cluster agent
  join:
    - 127.0.0.1:7947
```

### Environment Variables

| Environment Variable         | Description                              | Default        |
|------------------------------|------------------------------------------|----------------|
| HOSTNAME                     | Hostname to be set at startup.           | N/A            |
| RCOND_ADDR                   | Address to bind the HTTP server to.      | 0.0.0.0:8080   |
| RCOND_API_TOKEN              | API token to use for authentication.     | N/A            |
| RCOND_CLUSTER_ENABLED        | Enable the cluster agent.                | false          |
| RCOND_CLUSTER_NODE_NAME      | Name of the node in the cluster.         | rcond          |
| RCOND_CLUSTER_SECRET_KEY     | Secret key for the cluster agent.        | N/A            |
| RCOND_CLUSTER_ADVERTISE_ADDR | Advertise address for the cluster agent. | 0.0.0.0        |
| RCOND_CLUSTER_ADVERTISE_PORT | Advertise port for the cluster agent.    | 7946           |
| RCOND_CLUSTER_BIND_ADDR      | Bind address for the cluster agent.      | 0.0.0.0        |
| RCOND_CLUSTER_BIND_PORT      | Bind port for the cluster agent.         | 7946           |
| RCOND_CLUSTER_JOIN           | Join addresses for the cluster agent.    | 127.0.0.1:7947 |

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
| POST    | `/system/file`                      | Upload a file to the system             |
| POST    | `/system/restart`                   | Restart the system                      |
| POST    | `/system/shutdown`                  | Shutdown the system                     |
| GET     | `/cluster/members`                  | Get the cluster members                 |
| POST    | `/cluster/join`                     | Join cluster nodes                      |
| POST    | `/cluster/leave`                    | Leave the cluster                       |
| POST    | `/cluster/event`                    | Send a cluster event                    |


### Response Codes

- 200: Success
- 400: Bad request (invalid JSON payload)
- 405: Method not allowed
- 500: Internal server error

### Request/Response Format
All endpoints use JSON for request and response payloads.

## Cluster Events

Cluster events are used for broadcast messages to all nodes in the cluster. They are sent as HTTP POST requests to the `/cluster/event` endpoint.

The request body should be a JSON object with the following fields:

| Field    | Description                             | Optional  |
|----------|-----------------------------------------|-----------|
| `name`   | The name of the event                   | No        |
| `payload`| The payload of the event                | Yes       |

The response will be a JSON object with the following fields:

| Field    | Description                                                                      | Optional  |
|----------|----------------------------------------------------------------------------------|-----------|
| `status` | The status of the event. This is a string, either "success" or "error".          | No        |
| `error`  | If the status is "error", this field will contain a string describing the error. | Yes       |

Following events are implemented:

| Event Name | Description          | Payload |
|------------|----------------------|---------|
| restart    | Restart the cluster  | N/A     |
| shutdown   | Shutdown the cluster | N/A     |

## Examples

### Connect to a WiFi Access Point

This example will automatically connect to a WiFi access point with the given SSID and password on the interface "wlan0".

```bash
curl -X POST "http://rpi-test:8080/network/sta" \
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
curl -X POST "http://rpi-test:8080/network/ap" \
  -H "Content-Type: application/json" \
  -H "X-API-Token: 1234567890" \
  -d '{
    "interface": "wlan0",
    "ssid": "MyAccessPoint",
    "password": "StrongPassword",
    "autoconnect": true
  }'
```
### Restart the cluster

This example will restart all nodes in the cluster

```bash
curl -X POST "http://rpi-test:8080/cluster/event" \
  -H "accept: application/json" \
  -H "X-API-Token: 1234567890" \
  -d '{
    "name": "restart"
  }'
```

### Upload a file

This example will store Base64 encoded content to the target path.

```bash
curl -X 'POST' \
  'http://localhost:8080/system/file' \
  -H 'accept: application/json' \
  -H 'X-API-Token: 1234567890' \
  -H 'Content-Type: application/json' \
  -d '{
    "path": "/tmp/somefile",
    "content": "Zm9vCg=="
  }'
```