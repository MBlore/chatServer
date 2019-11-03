package main

import (
	"chatServer/server"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetOutput(os.Stdout)

	serv := server.NewTCPServer(
		onHandlePacket,
		onClientConnect,
		onClientDisconnect)

	serv.Listen(":80")
	serv.Run()
}

func onHandlePacket(s *server.TCPServer, client *server.Client, packet *server.Packet) {
	if !client.LoggedIn && packet.ID == PacketIDLogin {
		parts := strings.Split(string(*packet.Data), "\n")
		username := parts[0]
		password := parts[1]

		// TODO: Is user logged in already.
		// TODO: Validate the user from DB.
		// TODO: Yep, we definitly want our handlers in their own files.

		loggedIn := false

		if password == "pass" {
			loggedIn = true
			client.Username = username
			client.LoggedIn = true
			log.Printf("User '%v' logged in.", username)
		} else {
			log.Printf("User '%v' login denied.", username)
		}

		client.SendPacket(NewLoginResultPacket(loggedIn))
		return
	}

	// Logged in packet handlers...
	if client.LoggedIn {
		switch packet.ID {
		case PacketIDAudio:
			s.BroadcastPacket(&server.Packet{
				ID:   PacketIDAudio,
				Data: packet.Data,
			}, client)
		case PacketIDChat:
			if packet.Data != nil {
				s.BroadcastPacket(&server.Packet{
					ID:   PacketIDChat,
					Data: packet.Data,
				}, nil)
			}
		}
	}
}

func onClientConnect(s *server.TCPServer, c *server.Client, addr string) {
	log.Printf("Client connected from %v", addr)
}

func onClientDisconnect(s *server.TCPServer, c *server.Client, addr string) {
	log.Printf("Client disconnected from %v", addr)
}
