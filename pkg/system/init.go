package system

import (
	"log"

	"github.com/0x1d/rcond/pkg/config"
	"github.com/0x1d/rcond/pkg/network"
	"github.com/0x1d/rcond/pkg/util"
	"github.com/godbus/dbus/v5"
)

func Configure(appConfig *config.Config) error {
	log.Print("[INFO] Configure system")
	// configure hostname
	if err := network.SetHostname(appConfig.Hostname); err != nil {
		log.Printf("[ERROR] setting hostname failed: %s", err)
	}
	// configure network connections
	for _, connection := range appConfig.Network.Connections {
		err := util.WithConnection(func(conn *dbus.Conn) error {
			_, err := network.AddConnectionWithConfig(conn, &network.ConnectionConfig{
				Type:        connection.Type,
				UUID:        connection.UUID,
				ID:          connection.ID,
				AutoConnect: connection.AutoConnect,
				SSID:        connection.SSID,
				Mode:        connection.Mode,
				Band:        connection.Band,
				Channel:     connection.Channel,
				KeyMgmt:     connection.KeyMgmt,
				PSK:         connection.PSK,
				IPv4Method:  connection.IPv4Method,
				IPv6Method:  connection.IPv6Method,
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Printf("[ERROR] configuring connections failed: %s", err)
		}

	}
	log.Print("[INFO] System configured")
	return nil
}
