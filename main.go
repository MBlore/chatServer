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
	dbaccess.OpenConnection()
	defer dbaccess.CloseConnection()

	log.Println("Connected.")

	log.Println("Updating user statuses...")
	dbaccess.ResetUserStatuses()

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

		case builders.PacketIDConfirmContact:
			handlers.HandleConfirmContact(s, client, packet)

		case builders.PacketIDRejectContact:
			handlers.HandleRejectContact(s, client, packet)

		case builders.PacketIDImage:
			if packet.Data != nil {
				reader := bytes.NewReader(*packet.Data)

				userIDTo := utils.ReadInt32(reader)
				imageData := utils.ReadLenString(reader)

				if imageData != nil {
					// TODO: Validate this client is the senders friend.
					go s.BroadcastPacketToUserID(int(userIDTo), builders.NewImageFromPacket(client.UserID, imageData))
				}
			}
		case builders.PacketIDUserStatusChange:
			if packet.Data == nil {
				break
			}

			reader := bytes.NewReader(*packet.Data)
			statusID := utils.ReadInt32(reader)

			if statusID < 0 || statusID > 3 {
				break
			}

			err := dbaccess.SetStatus(client.UserID, int(statusID))

			if err != nil {
				log.Printf("Failed to status for user id '%v': %v", client.UserID, err)
			} else {
				packet := builders.NewUserStatusChangePacket(client.UserID, int(statusID))
				go utils.BroadcastPacketToContacts(s, client.UserID, packet)
			}
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
			err := dbaccess.LogoutUser(c.UserID)
			if err != nil {
				log.Print(err.Error())
				log.Printf("Failed to logout user id '%v' due to error.", c.UserID)
			}
		}()

		// Notify this contacts friends, if they are logged in, that this user went offline.
		statusPacket := builders.NewUserStatusChangePacket(c.UserID, dbaccess.StatusOffline)
		go utils.BroadcastPacketToContacts(s, c.UserID, statusPacket)
	}

	log.Printf("Client '%v' disconnected from %v (%v clients)", c.Username, addr, s.NumClients())
}
