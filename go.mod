module github.com/0x1d/rcond

go 1.23.4

replace github.com/0x1d/rcond/cmd => ./cmd

replace github.com/0x1d/rcond/pkg => ./pkg

require (
	github.com/godbus/dbus/v5 v5.1.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	golang.org/x/crypto v0.37.0
	gopkg.in/yaml.v3 v3.0.1
)

require golang.org/x/sys v0.32.0 // indirect
