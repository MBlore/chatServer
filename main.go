package main

import (
	"bytes"
	"chatServer/builders"
	"chatServer/dbaccess"
	"chatServer/handlers"
	"chatServer/server"
	"chatServer/utils"
	"io"
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Println("Connecting to database...")
	dbaccess.DBAccess.OpenConnection()
	defer dbaccess.DBAccess.CloseConnection()

	log.Println("Connected.")

	log.Println("Updating user statuses...")
	dbaccess.DBAccess.ResetUserStatuses()

	localTesting := true

	if !localTesting {
		go RunWebServer()
	}

	serv := server.NewTCPServer(
		onHandlePacket,
		onClientConnect,
		onClientDisconnect)

	log.Println("Starting chat server...")
	serv.Listen(":5035")
	serv.Run()
}

func onHandlePacket(s *server.TCPServer, client *server.Client, packet *server.Packet) {
	if !client.LoggedIn && packet.ID == builders.PacketIDLogin && packet.Data != nil {
		handlers.HandleLogin(s, client, packet)
		return
	}

	// Logged in packet handlers...
	if client.LoggedIn {
		switch packet.ID {
		case builders.PacketIDAudio:
			go s.BroadcastPacket(&server.Packet{
				ID:   builders.PacketIDAudio,
				Data: packet.Data,
			}, client)

		case builders.PacketIDChat:
			if packet.Data != nil {
				reader := bytes.NewReader(*packet.Data)

				userIDTo := utils.ReadInt32(reader)
				msg := utils.ReadLenString(reader)

				if msg != nil {
					// TODO: Validate this client is the senders friend.
					go s.BroadcastPacketToUserID(int(userIDTo), builders.NewChatFromPacket(client.UserID, *msg))
				}
			}

		case builders.PacketIDAction:
			if packet.Data != nil {
				reader := bytes.NewReader(*packet.Data)

				userIDTo := utils.ReadInt32(reader)
				action := utils.ReadLenString(reader)

				if action != nil {
					// TODO: Validate this client is the senders friend.
					go s.BroadcastPacketToUserID(int(userIDTo), builders.NewActionFromPacket(client.UserID, *action))
				}
			}

		case builders.PacketIDNudge:
			if packet.Data != nil {
				reader := bytes.NewReader(*packet.Data)

				userIDTo := utils.ReadInt32(reader)

				// TODO: Validate this client is the senders friend.

				go s.BroadcastPacketToUserID(int(userIDTo), builders.NewNudgeFromPacket(client.UserID))
			}

		case builders.PacketIDSetDisplayName:
			if packet.Data == nil {
				break
			}

			name := string(*packet.Data)
			client.DisplayName = &name

		case builders.PacketIDAddContact:
			handlers.HandleAddContact(s, client, packet)
		}
	}
}

func onClientConnect(s *server.TCPServer, c *server.Client, addr string) {
	log.Printf("Client connected from %v (%v clients)", addr, s.NumClients())
}

func onClientDisconnect(s *server.TCPServer, c *server.Client, addr string) {
	if c.LoggedIn && !s.IsClientMultiLogged(c.UserID) {
		// Only log out the user (set offline) if its the last client login instance.

		go func() {
			err := dbaccess.DBAccess.LogoutUser(c.UserID)
			if err != nil {
				log.Print(err.Error())
				log.Printf("Failed to logout user id '%v' due to error.", c.UserID)
			}
		}()

		// Notify this contacts friends, if they are logged in, that this user went offline.
		go func() {
			friends, err := dbaccess.DBAccess.GetFriends(c.UserID)
			if err != nil {
				log.Printf("Failed to get friends for user id '%v'.", c.UserID)
			} else {
				if friends != nil {
					statusPacket := builders.NewUserStatusChangePacket(c.UserID, dbaccess.StatusOffline)

					for _, f := range friends {
						go s.BroadcastPacketToUserID(int(f.ID), statusPacket)
					}
				}
			}
		}()
	}

	log.Printf("Client '%v' disconnected from %v (%v clients)", c.Username, addr, s.NumClients())
}
