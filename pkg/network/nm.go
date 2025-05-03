package network

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/godbus/dbus/v5"
)

const (
	ourUUID = "7d706027-727c-4d4c-a816-f0e1b99db8ab"
)

var (
	defaultSSID     = "PIAP"
	defaultPassword = "raspberry"
)

// withDbus executes the given function with a D-Bus system connection
// and handles any connection errors
func withDbus(fn func(*dbus.Conn) error) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Printf("Failed to connect to system bus: %v", err)
		return err
	}
	if err := fn(conn); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// ActivateConnection activates a NetworkManager connection profile.
// It takes a D-Bus connection, connection profile path, and device path as arguments.
// The function waits up to 10 seconds for the connection to become active.
// Returns an error if activation fails or times out.
func ActivateConnection(conn *dbus.Conn, connPath, devPath dbus.ObjectPath) error {
	nmObj := conn.Object(
		"org.freedesktop.NetworkManager",
		"/org/freedesktop/NetworkManager",
	)

	// Activate the connection
	var activePath dbus.ObjectPath
	err := nmObj.
		Call("org.freedesktop.NetworkManager.ActivateConnection", 0,
			connPath, devPath, dbus.ObjectPath("/")).
		Store(&activePath)
	if err != nil {
		return fmt.Errorf("ActivateConnection failed: %v", err)
	}

	// Wait until the connection is activated
	props := conn.Object(
		"org.freedesktop.NetworkManager",
		activePath,
	)

	start := time.Now()
	for time.Since(start) < 10*time.Second {
		var stateVar dbus.Variant
		err = props.
			Call("org.freedesktop.DBus.Properties.Get", 0,
				"org.freedesktop.NetworkManager.Connection.Active",
				"State").
			Store(&stateVar)
		if err != nil {
			return fmt.Errorf("Properties.Get(State) failed: %v", err)
		}
		if state, ok := stateVar.Value().(uint32); ok && state == 2 {
			log.Printf("Connection activated on connection path %v", connPath)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("failed to activate connection")
}

// DisconnectDevice disconnects a NetworkManager device, stopping any active connections.
// Takes a D-Bus connection and device path as arguments.
// Returns an error if the disconnect operation fails.
func DisconnectDevice(conn *dbus.Conn, devPath dbus.ObjectPath) error {
	devObj := conn.Object(
		"org.freedesktop.NetworkManager",
		devPath,
	)
	err := devObj.
		Call("org.freedesktop.NetworkManager.Device.Disconnect", 0).
		Err
	if err != nil {
		return fmt.Errorf("Device.Disconnect failed: %v", err)
	}
	fmt.Println("Access point stopped")
	return nil
}

// DeleteConnection removes a NetworkManager connection profile.
// Takes a D-Bus connection and connection profile path as arguments.
// Returns an error if the delete operation fails.
func DeleteConnection(conn *dbus.Conn, connPath dbus.ObjectPath) error {
	connObj := conn.Object(
		"org.freedesktop.NetworkManager",
		connPath,
	)
	err := connObj.
		Call("org.freedesktop.NetworkManager.Settings.Connection.Delete", 0).
		Err
	if err != nil {
		return fmt.Errorf("Connection.Delete failed: %v", err)
	}
	log.Printf("Connection removed: %v", connPath)
	return nil
}

// GetWifiConfig retrieves the WiFi SSID and password from environment variables or command line arguments.
// Takes an operation string ("up", "down", etc) to determine whether to check command line args.
// Returns the SSID and password strings.
// Priority order: command line args > environment variables > default values
func GetWifiConfig(op string) (string, string) {
	// Get SSID and password from args, env vars or use defaults
	ssid := defaultSSID
	password := defaultPassword

	// Check environment variables first
	if v, ok := os.LookupEnv("WIFI_SSID"); ok {
		ssid = v
	}
	if v, ok := os.LookupEnv("WIFI_PASSWORD"); ok {
		password = v
	}

	// Command line args override environment variables
	if op == "up" && len(os.Args) >= 5 {
		ssid = os.Args[3]
		password = os.Args[4]
	}

	return ssid, password
}

// GetConnectionPath looks up a NetworkManager connection profile by UUID.
// Takes a D-Bus connection and connection UUID string as arguments.
// Returns the D-Bus object path of the connection if found, or empty string if not found.
// Returns an error if the lookup operation fails.
func GetConnectionPath(conn *dbus.Conn, connUUID string) (dbus.ObjectPath, error) {
	// Get the Settings interface
	settingsObj := conn.Object(
		"org.freedesktop.NetworkManager",
		"/org/freedesktop/NetworkManager/Settings",
	)

	// List existing connections
	var paths []dbus.ObjectPath
	err := settingsObj.
		Call("org.freedesktop.NetworkManager.Settings.ListConnections", 0).
		Store(&paths)
	if err != nil {
		return "", fmt.Errorf("ListConnections failed: %v", err)
	}

	// Look up our connection by UUID
	var connPath dbus.ObjectPath
	for _, p := range paths {
		obj := conn.Object(
			"org.freedesktop.NetworkManager",
			p,
		)
		var cfg map[string]map[string]dbus.Variant
		err = obj.
			Call("org.freedesktop.NetworkManager.Settings.Connection.GetSettings", 0).
			Store(&cfg)
		if err != nil {
			continue
		}
		if v, ok := cfg["connection"]["uuid"].Value().(string); ok && v == connUUID {
			connPath = p
			break
		}
	}

	return connPath, nil
}

// AddConnection creates a new NetworkManager connection profile for a WiFi access point.
// Takes a D-Bus connection, SSID string, and password string as arguments.
// Returns the D-Bus object path of the new connection profile.
// Returns an error if the connection creation fails.
func AddConnection(conn *dbus.Conn, ssid string, password string) (dbus.ObjectPath, error) {
	settingsObj := conn.Object(
		"org.freedesktop.NetworkManager",
		"/org/freedesktop/NetworkManager/Settings",
	)

	settingsMap := map[string]map[string]dbus.Variant{
		"connection": {
			"type":        dbus.MakeVariant("802-11-wireless"),
			"uuid":        dbus.MakeVariant(ourUUID),
			"id":          dbus.MakeVariant(ssid),
			"autoconnect": dbus.MakeVariant(true),
		},
		"802-11-wireless": {
			"ssid":    dbus.MakeVariant([]byte(ssid)),
			"mode":    dbus.MakeVariant("ap"),
			"band":    dbus.MakeVariant("bg"),
			"channel": dbus.MakeVariant(uint32(1)),
		},
		"802-11-wireless-security": {
			"key-mgmt": dbus.MakeVariant("wpa-psk"),
			"psk":      dbus.MakeVariant(password),
		},
		"ipv4": {
			"method": dbus.MakeVariant("shared"),
		},
		"ipv6": {
			"method": dbus.MakeVariant("ignore"),
		},
	}

	var connPath dbus.ObjectPath
	err := settingsObj.
		Call("org.freedesktop.NetworkManager.Settings.AddConnection", 0, settingsMap).
		Store(&connPath)
	if err != nil {
		return "", fmt.Errorf("AddConnection failed: %v", err)
	}

	return connPath, nil
}

// GetDeviceByIpIface looks up a NetworkManager device by its interface name.
// Takes a D-Bus connection and interface name string as arguments.
// Returns the D-Bus object path of the device.
// Returns an error if the device lookup fails.
func GetDeviceByIpIface(conn *dbus.Conn, iface string) (dbus.ObjectPath, error) {
	// Get the NetworkManager interface
	nmObj := conn.Object(
		"org.freedesktop.NetworkManager",
		"/org/freedesktop/NetworkManager",
	)

	// Find the device by interface name
	var devPath dbus.ObjectPath
	err := nmObj.
		Call("org.freedesktop.NetworkManager.GetDeviceByIpIface", 0, iface).
		Store(&devPath)
	if err != nil {
		return "", fmt.Errorf("GetDeviceByIpIface(%s) failed: %v", iface, err)
	}

	return devPath, nil
}

// GetHostname returns the hostname of the current machine
func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("GetHostname failed: %v", err)
	}
	return hostname, nil
}

// SetHostname changes the static hostname via the system bus.
// newHost is your desired hostname, interactive=false skips any prompt.
func SetHostname(newHost string) error {
	return withDbus(func(conn *dbus.Conn) error {
		obj := conn.Object(
			"org.freedesktop.hostname1",
			dbus.ObjectPath("/org/freedesktop/hostname1"),
		)
		return obj.Call(
			"org.freedesktop.hostname1.SetStaticHostname",
			0,       // no special flags
			newHost, // the hostname you want
			false,   // interactive? (PolicyKit)
		).Err
	})
}

// Up creates and activates a WiFi access point connection.
// It takes the interface name, SSID, password and UUID as arguments.
// If a connection with the given UUID exists, it will be reused.
// Otherwise, a new connection will be created.
// The connection will be activated on the specified interface.
// Returns an error if any operation fails.
func Up(iface string, ssid string, password string, uuid string) error {
	return withDbus(func(conn *dbus.Conn) error {
		connPath, err := GetConnectionPath(conn, uuid)
		if err != nil {
			return err
		}

		if connPath == "" {
			connPath, err = AddConnection(conn, ssid, password)
			if err != nil {
				return err
			}
		}

		log.Printf("Getting device path for interface %s", iface)
		devPath, err := GetDeviceByIpIface(conn, iface)
		if err != nil {
			log.Printf("Failed to get device path for interface %s: %v", iface, err)
			return err
		}
		log.Printf("Got device path %s for interface %s", devPath, iface)

		if err := ActivateConnection(conn, connPath, devPath); err != nil {
			return err
		}

		return nil
	})
}

// Down deactivates a network connection on the specified interface.
// It takes the interface name as an argument.
// Returns an error if the device cannot be found or disconnected.
func Down(iface string) error {
	return withDbus(func(conn *dbus.Conn) error {
		devPath, err := GetDeviceByIpIface(conn, iface)
		if err != nil {
			return err
		}

		if err := DisconnectDevice(conn, devPath); err != nil {
			return err
		}

		return nil
	})
}

// Remove deletes a NetworkManager connection profile with the given UUID.
// If no connection with the UUID exists, it returns nil.
// Returns an error if the connection exists but cannot be deleted.
func Remove(uuid string) error {
	return withDbus(func(conn *dbus.Conn) error {
		connPath, err := GetConnectionPath(conn, uuid)
		if err != nil {
			return err
		}

		if connPath == "" {
			return nil
		}

		if err := DeleteConnection(conn, connPath); err != nil {
			return err
		}

		return nil
	})
}
