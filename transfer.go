package gows

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Handles the Message Frame
func handleMessage(buff []byte) ([]byte, error) {

	i := 0

	FIN_OPCODE_BYTE := buff[i]

	if !checkBit(FIN_OPCODE_BYTE, 7) {
		return nil, fmt.Errorf("fin not set to 1 (fragmeneted)")
	}
	i += 1

	MASK_LEN_BYTE := buff[i]

	if !checkBit(MASK_LEN_BYTE, 7) {
		return nil, fmt.Errorf("message frames has to be masked from client")
	}

	payload_len := int(MASK_LEN_BYTE & 0b01111111)

	if payload_len < 126 {
		fmt.Println("Got a Payload with len", payload_len)
		i += 1
	} else if payload_len == 126 {
		_bytes := buff[i+1 : i+3]
		x := binary.BigEndian.Uint16(_bytes)
		i += 3
		fmt.Println("Got a Payload with len", x)
	} else {
		_bytes := buff[i+1 : i+9]
		x := binary.BigEndian.Uint64(_bytes)
		fmt.Println("Got a Payload with len", x)
		i += 9
	}

	Mask := buff[i : i+4]

	i += 4

	data_buff := buff[i:]

	res := &bytes.Buffer{}

	for i, data := range data_buff {
		c := data ^ Mask[i%4]
		res.WriteByte(c)
	}

	return res.Bytes(), nil
}

func checkBit(b byte, position int) bool {

	return (b & (1 << position)) != 0

}
