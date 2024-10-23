package gows

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// Handles the Message Frame
func handleMessage(buff []byte) error {

	i := 0

	FIN_OPCODE_BYTE := fmt.Sprintf("%b", buff[i])
	i += 1

	if FIN_OPCODE_BYTE[0] != '1' {
		return fmt.Errorf("fin not set to 1 (fragmeneted)")
	}

	MASK_LEN_BYTE := fmt.Sprintf("%b", buff[i])

	if MASK_LEN_BYTE[0] != '1' {
		return fmt.Errorf("message frames has to be masked from client")
	}

	payload_len, _ := strconv.ParseInt(MASK_LEN_BYTE[i:], 2, 64)

	if payload_len < 126 {
		fmt.Println("Got a Payload with len", payload_len)
		i += 1
	} else if payload_len == 126 {
		x := uint16(buff[i+1])<<8 | uint16(buff[i+2])
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

	var decoded_data []byte

	for i, data := range data_buff {
		c := data ^ Mask[i%4]
		decoded_data = append(decoded_data, c)
	}

	fmt.Print("Recieved Data: ", string(decoded_data))

	return nil
}
