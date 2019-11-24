package handlers

import (
	"bytes"
	"chatServer/builders"
	"chatServer/dbaccess"
	"chatServer/server"
	"chatServer/utils"
)

// HandleRejectContact handles the receipt of a reject contact packet.
func HandleRejectContact(s *server.TCPServer, client *server.Client, packet *server.Packet) {
	if packet.Data == nil {
		return
	}

	reader := bytes.NewReader(*packet.Data)

	requestedUserID := int(utils.ReadInt32(reader))

	// Is this user pending for the client?
	rowID, err := dbaccess.GetPendingContact(requestedUserID, client.UserID)
	if err != nil {
		go client.SendPacket(builders.NewRejectContactResponsePacket(builders.RejectContactResultFailed, requestedUserID))
		return
	}

	if rowID == nil {
		go client.SendPacket(builders.NewRejectContactResponsePacket(builders.RejectContactResultNotPending, requestedUserID))
		return
	}

	// Do the reject.
	err = dbaccess.RejectContactRequest(requestedUserID, client.UserID)
	if err != nil {
		go client.SendPacket(builders.NewRejectContactResponsePacket(builders.RejectContactResultFailed, requestedUserID))
		return
	}

	go client.SendPacket(builders.NewRejectContactResponsePacket(builders.RejectContactResultSuccess, requestedUserID))
}
