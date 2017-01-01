package main

import (
	"encoding/binary"
	"fmt"
)

type packet struct {
	opcode     uint16 // two bytes
	blk_no     uint16 // two bytes
	dataLength uint16 // two bytes
	data       string
}

// Reference:
// http://stackoverflow.com/questions/26372227/go-conversion-between-struct-and-byte-array
func main() {
	str := "testing"

	pkt := packet{
		opcode:     3,
		blk_no:     1,
		dataLength: uint16(len(str)),
		data:       str,
	}

	// Convert struct to []byte
	var buf []byte = make([]byte, 50) // make sure the data string is less than 46 bytes

	offset := 0
	binary.BigEndian.PutUint16(buf[offset:], pkt.opcode)
	offset = offset + 2

	binary.BigEndian.PutUint16(buf[offset:], pkt.blk_no)
	offset = offset + 2

	binary.BigEndian.PutUint16(buf[offset:], pkt.dataLength)
	offset = offset + 2

	bytes_copied := copy(buf[offset:], pkt.data)

	fmt.Printf("Bytes copied: %v\n", bytes_copied)
	fmt.Printf("Buf: %v\n", buf)

	// Convert []byte back to struct
	_dataLength := binary.BigEndian.Uint16(buf[4:6])

	pkt2 := packet{
		opcode:     binary.BigEndian.Uint16(buf[0:2]), // two bytes
		blk_no:     binary.BigEndian.Uint16(buf[2:4]), // two bytes
		dataLength: _dataLength,                       // two bytes
		data:       string(buf[6 : 6+_dataLength]),
	}

	fmt.Printf("pkt2: %v\n", pkt2)
	fmt.Printf("pkt2: %#v\n", pkt2)
}
