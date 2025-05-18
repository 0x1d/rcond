module github.com/0x1d/rcond

go 1.23.4

replace github.com/0x1d/rcond/cmd => ./cmd

replace github.com/0x1d/rcond/pkg => ./pkg

require (
	github.com/godbus/dbus/v5 v5.1.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/hashicorp/logutils v1.0.0
	github.com/hashicorp/serf v0.10.2
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.37.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.4 // indirect
	github.com/hashicorp/go-msgpack/v2 v2.1.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-sockaddr v1.0.5 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/memberlist v0.5.2 // indirect
	github.com/miekg/dns v1.1.56 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529 // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/tools v0.14.0 // indirect
)
