package system

import (
	"log"

	"github.com/godbus/dbus/v5"
)

// WithDbus executes the given function with a D-Bus system connection
// and handles any connection errors
func WithDbus(fn func(*dbus.Conn) error) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("Failed to connect to system bus: %v", err)
		return err
	}
	if err := fn(conn); err != nil {
		log.Print(err)
		return err
	}
	conn.Close()
	return nil
}
