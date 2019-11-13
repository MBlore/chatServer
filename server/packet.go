package server

import "encoding/binary"

// Packet represents a received packet from a client connection.
type Packet struct {
	ID   int64
	Data *[]byte
}

// Converts a packet to a TLV byte buffer.
func (p *Packet) toBytes() *[]byte {
	packetIDBytes := make([]byte, 4)
	packetLengthBytes := make([]byte, 4)

	var dataLength uint32 = 0
	if p.Data != nil {
		dataLength = uint32(len(*p.Data))
	}

	binary.LittleEndian.PutUint32(packetIDBytes, uint32(p.ID))
	binary.LittleEndian.PutUint32(packetLengthBytes, dataLength)

	buffer := append(packetIDBytes, packetLengthBytes...)

	if p.Data != nil {
		buffer = append(buffer, *p.Data...)
	}

	return &buffer
}
