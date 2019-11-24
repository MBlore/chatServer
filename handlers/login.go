package handlers

import (
	"chatServer/builders"
	"chatServer/dbaccess"
	"chatServer/models"
	"chatServer/server"
	"chatServer/utils"
	"log"
	"strings"
)

// HandleLogin handles the receipt of a login packet.
func HandleLogin(s *server.TCPServer, client *server.Client, packet *server.Packet) {
	parts := strings.Split(string(*packet.Data), "\n")
	username := parts[0]
	password := parts[1]

	user, err := dbaccess.GetUserByUsername(username)
	if err != nil {
		panic(err.Error())
	}

	loggedIn := false
	var friends []models.FriendModel
	var pendingContacts []models.PendingContactModel

	if user != nil && utils.ComparePasswordHashes(password, user.Password) {
		err := dbaccess.LoginUser(user.ID)
		if err != nil {
			log.Print(err.Error())
			log.Printf("User '%v' login failed due to error.", username)
		} else {
			loggedIn = true
			client.Username = user.Username
			client.DisplayName = user.DisplayName
			client.LoggedIn = true
			client.UserID = user.ID
			client.Status = 1 // Default status set from the LoginUser proc.
			client.ImageURL = user.ImageURL
			client.StatusText = user.StatusText

			log.Printf("User '%v' logged in.", username)

			f, err := dbaccess.GetFriends(user.ID)
			if err != nil {
				log.Printf("Failed to get friends list for user ID '%v'.", user.ID)
			} else {
				friends = f
			}

			pc, err := dbaccess.GetUserPendingContacts(user.ID)
			if err != nil {
				log.Printf("Failed to get pending contacts list for user ID '%v'.", user.ID)
			} else {
				pendingContacts = pc
			}
		}
	} else {
		log.Printf("User '%v' login denied.", username)
	}

	if loggedIn {
		go client.SendPacket(builders.NewLoginResultPacket(
			loggedIn, client.UserID, client.DisplayName, user.StatusText, friends, pendingContacts))

		// Notify this contacts friends, if they are logged in, that this user came online.
		go func() {
			if friends != nil {
				statusPacket := builders.NewUserStatusChangePacket(user.ID, dbaccess.StatusOnline)

				for _, f := range friends {
					go s.BroadcastPacketToUserID(f.ID, statusPacket)
				}
			}
		}()
	} else {
		go client.SendPacket(builders.NewLoginResultPacket(false, 0, nil, nil, nil, nil))
	}
}
