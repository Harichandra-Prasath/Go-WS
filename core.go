package gows

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
)

var hasher = sha1.New()

// Handles the Initial Handshake parsing the headers
func handleHandShake(conn *net.TCPConn) error {

	// Read the Initial Handshake
	buff := make([]byte, 1024)
	buffer := &bytes.Buffer{}

	n, err := conn.Read(buff)
	if err != nil {
		return fmt.Errorf("reading handshake: %s", err)
	}

	headerMap := make(map[string]string)

	err = parseHandShake(string(buff[:n]), &headerMap)
	if err != nil {
		writeStatusLine("400", "Bad Request", buffer)
		buffer.WriteString("\r\n")
		conn.Write(buffer.Bytes())
		return fmt.Errorf("error in parseHandShake: %s", err)
	}

	fmt.Println("HandShaked Parsed and found to be Valid")

	respData := writeHandShake(&headerMap, buffer)

	_, err = conn.Write(respData)
	if err != nil {
		return fmt.Errorf("sending server handshake: %s", err)
	}

	fmt.Println("Wrote Server Handshake back")

	return nil
}

// Parser for Initial HandShake
func parseHandShake(data string, _headerMap *map[string]string) error {

	headerMap := *(_headerMap)

	// Get the headers
	_headers := strings.Split(data, "\r\n")

	l := len(_headers)

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

	if err := checkHeader("Upgrade", "websocket", _headerMap); err != nil {
		return err
	}
	if err := checkHeader("Connection", "Upgrade", _headerMap); err != nil {
		return err
	}
	if err := checkHeader("Sec-WebSocket-Version", "13", _headerMap); err != nil {
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

// Checks against expected value for the given header
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

// Sends server Handshake to the client
func writeHandShake(_headerMap *map[string]string, buff *bytes.Buffer) []byte {

	headerMap := *(_headerMap)

	writeStatusLine("101", "Switching Protocols", buff)

	// Prepare the Sec-WebSocket-Accept
	key := headerMap["Sec-WebSocket-Key"]
	key += ACCEPT_STRING

	hasher.Write([]byte(key))

	Accept_key := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Prepare the response headers

	respHeaderMap := make(map[string]string)

	respHeaderMap["Upgrade"] = "websocket"
	respHeaderMap["Connection"] = "Upgrade"
	respHeaderMap["Sec-WebSocket-Accept"] = Accept_key

	writeHeaders(&respHeaderMap, buff)
	buff.WriteString("\r\n")

	return buff.Bytes()

}

// Write the statusLine for the responses
func writeStatusLine(status string, text string, buff *bytes.Buffer) {

	buff.WriteString(fmt.Sprintf("HTTP/1.1 %s", status))
	buff.WriteByte(' ')
	buff.WriteString(text)
	buff.WriteString("\r\n")

}

// For given Headers, Write it to the buffer for final conn write
func writeHeaders(_headerMap *map[string]string, buff *bytes.Buffer) {

	headerMap := *(_headerMap)

	for key, value := range headerMap {
		buff.WriteString(fmt.Sprintf("%s:", key))
		buff.WriteByte(' ')
		buff.WriteString(value)
		buff.WriteString("\r\n")
	}

}
