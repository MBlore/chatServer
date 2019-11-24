package handlers

import (
	"bytes"
	"chatServer/builders"
	"chatServer/dbaccess"
	"chatServer/server"
	"chatServer/utils"
)

// HandleAddContact handles the receipt of a add contact packet.
func HandleAddContact(s *server.TCPServer, client *server.Client, packet *server.Packet) {
	if packet.Data != nil {
		reader := bytes.NewReader(*packet.Data)

		usernameToAdd := utils.ReadLenString(reader)
		message := utils.ReadLenString(reader)

		if usernameToAdd == nil {
			go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultUserNotFound))
			return
		}

		// Verify user exists.
		user, err := dbaccess.GetUserByUsername(*usernameToAdd)
		if err != nil {
			go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultFailed))
			return
		}

		if user == nil {
			go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultUserNotFound))
			return
		}

		// Can't add yourself...
		if user.ID == client.UserID {
			go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultUserNotFound))
			return
		}

		// Verify user is not already a contact.
		rowID, err := dbaccess.GetUserContactByContactUserID(client.UserID, user.ID)
		if err != nil {
			go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultFailed))
			return
		}

		if rowID != nil {
			go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultUserAlreadyContact))
			return
		}

		// Verify user is not already a pending contact.
		rowID, err = dbaccess.GetPendingContact(client.UserID, user.ID)
		if err != nil {
			go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultFailed))
			return
		}

		if rowID != nil {
			go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultUserAlreadyPending))
			return
		}

		// Validations passed - add the contact as pending.
		err = dbaccess.AddPendingContact(client.UserID, user.ID, message)
		if err != nil {
			go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultFailed))
			return
		}

		go client.SendPacket(builders.NewAddContactResponsePacket(builders.AddContactResultSuccess))

		// Notify contact they have a pending request.
		go s.BroadcastPacketToUserID(user.ID, builders.NewNotifyAddRequestPacket(client.UserID, &client.Username, client.DisplayName, message))
	}
}
