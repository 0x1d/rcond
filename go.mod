module github.com/0x1d/rcond

go 1.23.4

replace github.com/0x1d/rcond/cmd => ./cmd

replace github.com/0x1d/rcond/pkg => ./pkg

require (
	github.com/godbus/dbus/v5 v5.1.0
	github.com/gorilla/mux v1.8.1
)
