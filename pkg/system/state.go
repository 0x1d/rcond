package system

import (
	"log"

	"github.com/0x1d/rcond/pkg/util"
	"github.com/godbus/dbus/v5"
)

// Restart restarts the system.
func Restart() error {
	return util.WithConnection(func(conn *dbus.Conn) error {
		obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
		log.Println("Rebooting system...")
		call := obj.Call("org.freedesktop.systemd1.Manager.Reboot", 0)
		if call.Err != nil {
			return call.Err
		}
		return nil
	})
}

// Shutdown shuts down the system.
func Shutdown() error {
	return util.WithConnection(func(conn *dbus.Conn) error {
		obj := conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
		log.Println("Shutting down system...")
		call := obj.Call("org.freedesktop.systemd1.Manager.PowerOff", 0)
		if call.Err != nil {
			return call.Err
		}
		return nil
	})
}
