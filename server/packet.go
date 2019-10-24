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

	binary.LittleEndian.PutUint32(packetIDBytes, uint32(p.ID))
	binary.LittleEndian.PutUint32(packetLengthBytes, uint32(len(*p.Data)))

	buffer := append(packetIDBytes, packetLengthBytes...)
	buffer = append(buffer, *p.Data...)

	return &buffer
}
