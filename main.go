package main

import (
	"bytes"
	"chatServer/server"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Println("Connecting to database...")
	DBAccess.OpenConnection()
	defer DBAccess.CloseConnection()

	log.Println("Connected.")

	log.Println("Updating user statuses...")
	DBAccess.ResetUserStatuses()

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
	if !client.LoggedIn && packet.ID == PacketIDLogin && packet.Data != nil {
		parts := strings.Split(string(*packet.Data), "\n")
		username := parts[0]
		password := parts[1]

		user, err := DBAccess.GetUserByUsername(username)
		if err != nil {
			panic(err.Error())
		}

		loggedIn := false
		var friends []FriendModel

		if user != nil && ComparePasswordHashes(password, user.Password) {
			err := DBAccess.LoginUser(user.ID)
			if err != nil {
				log.Print(err.Error())
				log.Printf("User '%v' login failed due to error.", username)
			} else {
				loggedIn = true
				client.Username = user.Username
				client.DisplayName = user.DisplayName
				client.LoggedIn = true
				client.UserID = user.ID
				log.Printf("User '%v' logged in.", username)

				f, err := DBAccess.GetFriends(user.ID)
				if err != nil {
					log.Printf("Failed to get friends list for user ID '%v'.", user.ID)
				} else {
					friends = f
				}
			}
		} else {
			log.Printf("User '%v' login denied.", username)
		}

		if loggedIn {
			go client.SendPacket(NewLoginResultPacket(loggedIn, client.UserID, client.DisplayName, friends))

			// Notify this contacts friends, if they are logged in, that this user came online.
			go func() {
				if friends != nil {
					statusPacket := NewUserStatusChangePacket(user.ID, StatusOnline)

					for _, f := range friends {
						go s.BroadcastPacketToUserID(f.ID, statusPacket)
					}
				}
			}()
		} else {
			go client.SendPacket(NewLoginResultPacket(false, 0, nil, nil))
		}
		return
	}

	// Logged in packet handlers...
	if client.LoggedIn {
		switch packet.ID {
		case PacketIDAudio:
			go s.BroadcastPacket(&server.Packet{
				ID:   PacketIDAudio,
				Data: packet.Data,
			}, client)

		case PacketIDChat:
			if packet.Data != nil {
				reader := bytes.NewReader(*packet.Data)

				userIDTo := readInt32(reader)
				msg := readLenString(reader)

				if msg != nil {
					// TODO: Validate this client is the senders friend.
					go s.BroadcastPacketToUserID(int(userIDTo), NewChatFromPacket(client.UserID, *msg))
				}
			}

		case PacketIDAction:
			if packet.Data != nil {
				reader := bytes.NewReader(*packet.Data)

				userIDTo := readInt32(reader)
				action := readLenString(reader)

				if action != nil {
					// TODO: Validate this client is the senders friend.
					go s.BroadcastPacketToUserID(int(userIDTo), NewActionFromPacket(client.UserID, *action))
				}
			}

		case PacketIDNudge:
			if packet.Data != nil {
				reader := bytes.NewReader(*packet.Data)

				userIDTo := readInt32(reader)

				// TODO: Validate this client is the senders friend.

				go s.BroadcastPacketToUserID(int(userIDTo), NewNudgeFromPacket(client.UserID))
			}

		case PacketIDSetDisplayName:
			if packet.Data == nil {
				break
			}

			name := string(*packet.Data)
			client.DisplayName = &name

		case PacketIDAddContact:
			if packet.Data != nil {
				reader := bytes.NewReader(*packet.Data)

				usernameToAdd := readLenString(reader)
				message := readLenString(reader)

				if usernameToAdd == nil {
					client.SendPacket(NewAddContactResponsePacket(AddContactResultUserNotFound))
					break
				}

				// Verify user exists.
				user, err := DBAccess.GetUserByUsername(*usernameToAdd)
				if err != nil {
					client.SendPacket(NewAddContactResponsePacket(AddContactResultFailed))
					break
				}

				if user == nil {
					client.SendPacket(NewAddContactResponsePacket(AddContactResultUserNotFound))
					break
				}

				// Verify user is not already a contact.
				rowID, err := DBAccess.GetUserContactByContactUserID(client.UserID, user.ID)
				if err != nil {
					client.SendPacket(NewAddContactResponsePacket(AddContactResultFailed))
					break
				}

				if rowID != nil {
					client.SendPacket(NewAddContactResponsePacket(AddContactResultUserAlreadyContact))
					break
				}

				// Verify user is not already a pending contact.
				rowID, err = DBAccess.GetPendingContact(client.UserID, user.ID)
				if err != nil {
					client.SendPacket(NewAddContactResponsePacket(AddContactResultFailed))
					break
				}

				if rowID != nil {
					client.SendPacket(NewAddContactResponsePacket(AddContactResultUserAlreadyPending))
					break
				}

				// Validations passed - add the contact as pending.
				err = DBAccess.AddPendingContact(client.UserID, user.ID, message)
				if err != nil {
					client.SendPacket(NewAddContactResponsePacket(AddContactResultFailed))
					break
				}

				client.SendPacket(NewAddContactResponsePacket(AddContactResultSuccess))
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
			err := DBAccess.LogoutUser(c.UserID)
			if err != nil {
				log.Print(err.Error())
				log.Printf("Failed to logout user id '%v' due to error.", c.UserID)
			}
		}()

		// Notify this contacts friends, if they are logged in, that this user went offline.
		go func() {
			friends, err := DBAccess.GetFriends(c.UserID)
			if err != nil {
				log.Printf("Failed to get friends for user id '%v'.", c.UserID)
			} else {
				if friends != nil {
					statusPacket := NewUserStatusChangePacket(c.UserID, StatusOffline)

					for _, f := range friends {
						go s.BroadcastPacketToUserID(int(f.ID), statusPacket)
					}
				}
			}
		}()
	}

	log.Printf("Client '%v' disconnected from %v (%v clients)", c.Username, addr, s.NumClients())
}
