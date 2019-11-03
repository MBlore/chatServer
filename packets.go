package main

import "chatServer/server"

// Packet types supported in our client/server protocol.
const (
	PacketIDLogin       = 0
	PacketIDLoginResult = 1
	PacketIDChat        = 2
	PacketIDNudge       = 3
	PacketIDAudio       = 4
	PacketIDPing        = 5
)

// NewLoginResultPacket creates a new login result packet.
func NewLoginResultPacket(success bool) *server.Packet {

	data := make([]byte, 1)

	if success {
		data[0] = 1
	} else {
		data[0] = 0
	}

	packet := server.Packet{
		ID:   PacketIDLoginResult,
		Data: &data,
	}

	return &packet
}
