package handlers

import (
	"bytes"
	"chatServer/builders"
	"chatServer/dbaccess"
	"chatServer/server"
	"chatServer/utils"
)

// HandleConfirmContact handles the receipt of a confirm contact packet.
func HandleConfirmContact(s *server.TCPServer, client *server.Client, packet *server.Packet) {
	if packet.Data == nil {
		return
	}

	reader := bytes.NewReader(*packet.Data)

	requestedUserID := int(utils.ReadInt32(reader))

	// Is this user pending for the client?
	rowID, err := dbaccess.GetPendingContact(requestedUserID, client.UserID)
	if err != nil {
		go client.SendPacket(builders.NewConfirmContactResponsePacket(builders.ConfirmContactResultFailed, requestedUserID, nil))
		return
	}

	if rowID == nil {
		go client.SendPacket(builders.NewConfirmContactResponsePacket(builders.ConfirmContactResultNotPending, requestedUserID, nil))
		return
	}

	// Do the confirm.
	err = dbaccess.ConfirmContactRequest(requestedUserID, client.UserID)
	if err != nil {
		go client.SendPacket(builders.NewConfirmContactResponsePacket(builders.ConfirmContactResultFailed, requestedUserID, nil))
		return
	}

	// Fetch the details of the user so the client can add them to the contacts.
	user, err := dbaccess.GetUserByID(requestedUserID)
	if err != nil || user == nil {
		go client.SendPacket(builders.NewConfirmContactResponsePacket(builders.ConfirmContactResultFailed, requestedUserID, nil))
		return
	}

	go client.SendPacket(builders.NewConfirmContactResponsePacket(builders.ConfirmContactResultSuccess, requestedUserID, user))

	// Now we need to tell the requester that this client accepted.
	go s.BroadcastPacketToUserID(requestedUserID, builders.NewAddContactAcceptedPacket(
		client.UserID, &client.Username, client.DisplayName, client.Status, client.ImageURL, client.StatusText))
}
