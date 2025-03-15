package utils

import (
	"net"

	"broadcast-server/config"
)

func IsPortTaken(port string) bool {
	addr, err := net.Listen(config.DefaultServerProtocol, ":"+port)
	if err != nil {
		return true // Port is taken
	}
	defer addr.Close() // Now it's safe, because Listen() succeeded
	return false       // Port is free
}
