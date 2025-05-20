package network

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/0x1d/rcond/pkg/util"
	"github.com/godbus/dbus/v5"
	"github.com/google/uuid"
)

var (
	defaultSSID     = "PIAP"
	defaultPassword = "raspberry"
)

// ConnectionConfig holds the configuration for a NetworkManager connection
type ConnectionConfig struct {
	Type        string
	UUID        string
	ID          string
	AutoConnect bool
	SSID        string
	Mode        string
	Band        string
	Channel     uint32
	KeyMgmt     string
	PSK         string
	IPv4Method  string
	IPv6Method  string
}

func DefaultSTAConfig(uuid uuid.UUID, ssid string, password string, autoconnect bool) *ConnectionConfig {
	return &ConnectionConfig{
		Type:        "802-11-wireless",
		UUID:        uuid.String(),
		ID:          ssid,
		AutoConnect: autoconnect,
		SSID:        ssid,
		Mode:        "infrastructure",
		KeyMgmt:     "wpa-psk",
		PSK:         password,
		IPv4Method:  "auto",
		IPv6Method:  "ignore",
	}
}

// DefaultAPConfig returns a default access point configuration
func DefaultAPConfig(uuid uuid.UUID, ssid string, password string, autoconnect bool) *ConnectionConfig {
	return &ConnectionConfig{
		Type:        "802-11-wireless",
		UUID:        uuid.String(),
		ID:          ssid,
		AutoConnect: autoconnect,
		SSID:        ssid,
		Mode:        "ap",
		Band:        "bg",
		Channel:     1,
		KeyMgmt:     "wpa-psk",
		PSK:         password,
		IPv4Method:  "shared",
		IPv6Method:  "ignore",
	}
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
// Takes a D-Bus connection, UUID string, SSID string, and password string as arguments.
// Returns the D-Bus object path of the new connection profile.
// Returns an error if the connection creation fails.
func AddAccessPointConnection(conn *dbus.Conn, uuid uuid.UUID, ssid string, password string, autoconnect bool) (dbus.ObjectPath, error) {
	return AddConnectionWithConfig(conn, DefaultAPConfig(uuid, ssid, password, autoconnect))
}

// AddStationConnection creates a new NetworkManager connection profile for a WiFi station (client).
// Takes a D-Bus connection, UUID string, SSID string, and password string as arguments.
// Returns the D-Bus object path of the new connection profile.
// Returns an error if the connection creation fails.
func AddStationConnection(conn *dbus.Conn, uuid uuid.UUID, ssid string, password string, autoconnect bool) (dbus.ObjectPath, error) {
	return AddConnectionWithConfig(conn, DefaultSTAConfig(uuid, ssid, password, autoconnect))
}

// AddConnectionWithConfig creates a new NetworkManager connection profile with the given configuration.
// Takes a D-Bus connection and ConnectionConfig struct as arguments.
// Returns the D-Bus object path of the new connection profile.
// Returns an error if the connection creation fails.
func AddConnectionWithConfig(conn *dbus.Conn, cfg *ConnectionConfig) (dbus.ObjectPath, error) {

	// check of connection already exists and return existing connection path
	if existingObjectPath, err := GetConnectionPath(conn, cfg.UUID); err == nil {
		return existingObjectPath, nil
	}

	settingsObj := conn.Object(
		"org.freedesktop.NetworkManager",
		"/org/freedesktop/NetworkManager/Settings",
	)

	settingsMap := map[string]map[string]dbus.Variant{
		"connection": {
			"type":        dbus.MakeVariant(cfg.Type),
			"uuid":        dbus.MakeVariant(cfg.UUID),
			"id":          dbus.MakeVariant(cfg.ID),
			"autoconnect": dbus.MakeVariant(cfg.AutoConnect),
		},
		"ipv4": {
			"method": dbus.MakeVariant(cfg.IPv4Method),
		},
		"ipv6": {
			"method": dbus.MakeVariant(cfg.IPv6Method),
		},
	}

	// configure wireless
	if cfg.Type == "802-11-wireless" {
		var wirelessMap map[string]dbus.Variant
		if cfg.Mode == "ap" {
			// Access Point mode
			wirelessMap = map[string]dbus.Variant{
				"ssid":    dbus.MakeVariant([]byte(cfg.SSID)),
				"mode":    dbus.MakeVariant(cfg.Mode),
				"band":    dbus.MakeVariant(cfg.Band),
				"channel": dbus.MakeVariant(cfg.Channel),
			}
		} else {
			// Station / Infrastructure mode
			wirelessMap = map[string]dbus.Variant{
				"ssid": dbus.MakeVariant([]byte(cfg.SSID)),
				"mode": dbus.MakeVariant(cfg.Mode),
			}
		}
		settingsMap["802-11-wireless"] = wirelessMap
		settingsMap["802-11-wireless-security"] = map[string]dbus.Variant{
			"key-mgmt": dbus.MakeVariant(cfg.KeyMgmt),
			"psk":      dbus.MakeVariant(cfg.PSK),
		}
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
	return util.WithConnection(func(conn *dbus.Conn) error {
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

// ConfigureSTA connects to a WiFi access point with the specified settings.
// It takes the interface name, SSID and password as arguments.
// A new connection with a generated UUID will be created.
// Returns the UUID of the created connection and any error that occurred.
func ConfigureSTA(iface string, ssid string, password string, autoconnect bool) (string, error) {
	uuid := uuid.New()

	err := util.WithConnection(func(conn *dbus.Conn) error {
		_, err := AddStationConnection(conn, uuid, ssid, password, autoconnect)
		if err != nil {
			return fmt.Errorf("failed to create station connection: %v", err)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// ConfigureAP creates a WiFi access point connection with the specified settings.
// It takes the interface name, SSID and password as arguments.
// A new connection with a generated UUID will be created.
// Returns the UUID of the created connection and any error that occurred.
func ConfigureAP(iface string, ssid string, password string, autoconnect bool) (string, error) {
	uuid := uuid.New()

	err := util.WithConnection(func(conn *dbus.Conn) error {
		_, err := AddAccessPointConnection(conn, uuid, ssid, password, autoconnect)
		if err != nil {
			return fmt.Errorf("failed to create access point connection: %v", err)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// Up activates a connection.
// It takes the interface name and UUID as arguments.
// The connection with the given UUID must exist.
// The connection will be activated on the specified interface.
// Returns an error if any operation fails.
func Up(iface string, uuid string) error {
	return util.WithConnection(func(conn *dbus.Conn) error {
		connPath, err := GetConnectionPath(conn, uuid)
		if err != nil {
			return err
		}

		if connPath == "" {
			return fmt.Errorf("connection with UUID %s not found", uuid)
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
	return util.WithConnection(func(conn *dbus.Conn) error {
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
	return util.WithConnection(func(conn *dbus.Conn) error {
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
