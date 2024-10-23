package gows

import "bytes"

func invalidHandshake() []byte {

	buff := &bytes.Buffer{}

	buff.WriteString("HTTP/1.1 400")
	buff.WriteByte(' ')
	buff.WriteString("Bad Request")
	buff.WriteString("\r\n")

	return buff.Bytes()
}
