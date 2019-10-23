package main

import (
	"chatServer/server"
	"log"
	"net"
	"os"
)

// Packet types supported in our client/server protocol.
const (
	packetAudio = 1
	packetChat  = 2
)

func main() {
	log.SetOutput(os.Stdout)

	serv := server.NewTCPServer(onHandlePacket, onClientConnect)
	serv.Listen(":80")
	serv.Run()
}

func onHandlePacket(s *server.TCPServer, c net.Conn, p server.Packet) {
	// TODO: Make packet handler interface and have handlers in their own files, maybe?

	switch p.ID {
	case packetAudio:
		// Reply the packet to other connected clients.
		s.BroadcastPacket(server.Packet{
			ID:   packetAudio,
			Data: p.Data,
		}, c)

	case packetChat:
		// Send the chat to the other clients.
		s.BroadcastPacket(server.Packet{
			ID:   packetChat,
			Data: p.Data,
		}, c)
	}
}

func onClientConnect(s *server.TCPServer, c net.Conn, addr string) {
	// When clients connect, we may want to perform a handshake such as sending info about the server.
	log.Printf("Client connected from %v", addr)
}