package gows

import (
	"fmt"
	"net"
)

// Handles the Initial Handshake parsing the headers
func handleHandShake(conn *net.TCPConn) error {

	// Read the Initial Handshake
	buff := make([]byte, 1024)

	n, err := conn.Read(buff)
	if err != nil {
		return fmt.Errorf("reading handshake: %s", err)
	}

	conn.Write(buff[:n])

	return nil
}
