package util

import (
	"log"

	"github.com/godbus/dbus/v5"
)

// WithConnection executes the given function with a D-Bus system connection
// and handles any connection errors
func WithConnection(fn func(*dbus.Conn) error) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("[ERROR] Failed to connect to system bus: %v", err)
		return err
	}
	if err := fn(conn); err != nil {
		log.Printf("[ERROR] Failed to execute D-Bus function: %s", err)
		return err
	}
	conn.Close()
	return nil
}
