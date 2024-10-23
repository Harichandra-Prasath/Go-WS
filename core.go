package gows

import (
	"encoding/base64"
	"fmt"
	"net"
	"strings"
)

// Handles the Initial Handshake parsing the headers
func handleHandShake(conn *net.TCPConn) error {

	// Read the Initial Handshake
	buff := make([]byte, 1024)

	n, err := conn.Read(buff)
	if err != nil {
		return fmt.Errorf("reading handshake: %s", err)
	}

	err = parseHandShake(string(buff[:n]))
	if err != nil {
		return fmt.Errorf("parsing handshake: %s", err)
	}

	fmt.Println("HandShaked Parsed and found to be Valid")

	return nil
}

// Parse the HandShake
func parseHandShake(data string) error {

	// Get the headers
	_headers := strings.Split(data, "\r\n")

	l := len(_headers)

	headerMap := make(map[string]string)

	request_line := _headers[0]

	if !strings.Contains(request_line, "GET") {
		return fmt.Errorf("invalid method for handshake")
	}

	for _, header := range _headers[1 : l-2] {

		parts := strings.Split(header, ":")

		cleaned_key := strings.TrimSpace(parts[0])
		cleaned_value := strings.TrimSpace(parts[1])
		headerMap[cleaned_key] = cleaned_value

	}

	if err := checkHeader("Upgrade", "websocket", &headerMap); err != nil {
		return err
	}
	if err := checkHeader("Connection", "Upgrade", &headerMap); err != nil {
		return err
	}
	if err := checkHeader("Sec-WebSocket-Version", "13", &headerMap); err != nil {
		return err
	}

	// check for the length of the decoded key

	key_string, ok := headerMap["Sec-WebSocket-Key"]
	if !ok {
		return fmt.Errorf("no Sec-WebSocket-Key header present")
	}

	key_bytes, _ := base64.StdEncoding.DecodeString(key_string)

	if len(key_bytes) != 16 {
		return fmt.Errorf("invalid Sec-WebSocket-Key header")
	}

	return nil
}

func checkHeader(key string, expected string, _headerMap *map[string]string) error {

	headerMap := *(_headerMap)

	if value, ok := headerMap[key]; ok {

		if !strings.Contains(value, expected) && value != expected {
			return fmt.Errorf("invlaid %s header", key)
		}
	} else {
		return fmt.Errorf("no %s header  present", key)
	}

	return nil

}
