package builders

import (
	"bytes"
	"chatServer/models"
	"chatServer/server"
	"encoding/binary"
)

// Packet types supported in our client/server protocol.
const (
	PacketIDLogin                  = 0
	PacketIDLoginResult            = 1
	PacketIDChat                   = 2
	PacketIDNudge                  = 3
	PacketIDAudio                  = 4
	PacketIDPing                   = 5
	PacketIDAction                 = 6
	PacketIDActionFrom             = 7
	PacketIDSetDisplayName         = 8
	PacketIDNudgeFrom              = 9
	PacketIDUserStatusChange       = 10
	PacketIDChatFrom               = 11
	PacketIDAddContact             = 12
	PacketIDAddContactResponse     = 13
	PacketIDNotifyAddRequest       = 14
	PacketIDConfirmContact         = 15
	PacketIDConfirmContactResponse = 16
	PacketIDRejectContact          = 17
	PacketIDRejectContactResponse  = 18
	PacketIDAddContactAccepted     = 19
	PacketIDImage                  = 20
	PacketIDImageFrom              = 21
)

// Result codes for an add contact request.
const (
	AddContactResultFailed             = 0
	AddContactResultSuccess            = 1
	AddContactResultUserNotFound       = 2
	AddContactResultUserAlreadyContact = 3
	AddContactResultUserAlreadyPending = 4
)

// Result codes for a confirm contact request.
const (
	ConfirmContactResultSuccess    = 0
	ConfirmContactResultNotPending = 1
	ConfirmContactResultFailed     = 2
)

// Result codes for a reject contact request.
const (
	RejectContactResultSuccess    = 0
	RejectContactResultNotPending = 1
	RejectContactResultFailed     = 2
)

// Writes a string to the specified buffer prefixed by its length. If the string is nill, a length of 0 is written and no string data.
func writeString(buf *bytes.Buffer, str *string) {
	if str != nil {
		length := int32(len(*str))
		binary.Write(buf, binary.LittleEndian, int32(length))
		buf.Write([]byte(*str))
	} else {
		binary.Write(buf, binary.LittleEndian, int32(0))
	}
}

func writeBool(buf *bytes.Buffer, b bool) {
	binary.Write(buf, binary.LittleEndian, b)
}

func writeInt32(buf *bytes.Buffer, val int) {
	binary.Write(buf, binary.LittleEndian, int32(val))
}

// NewLoginResultPacket creates a new login result packet.
func NewLoginResultPacket(
	success bool,
	userID int,
	displayName *string,
	statusText *string,
	friends []models.FriendModel,
	pendingContacts []models.PendingContactModel) *server.Packet {
	/*
		Success (byte)
		UserId (int32)
		DisplayNameLength (int32)
		DisplayName (string)
		StatusTextLen (int32)
		StatusText

		FriendsCount (int32)

		(For each friend...)
		ID (int32)
		UsernameLength (int32)
		Username (string)
		DisplayNameLength (int32)
		DisplayName (string)
		Status (int32)
		ImageURLLength (int32)
		ImageURL (string)
		StatusTextLength (int32)
		StatusText (string)

		PendingContactsCount (int32)

		(For each pending contant)
		ID (int32)
		UsernameLength (int32)
		Username (string)
		DisplayNameLength (int32)
		DisplayName (string)
		ImageURLLength (int32)
		ImageURL (string)
		MessageLength (int32)
		Message (string)
	*/

	buf := new(bytes.Buffer)

	writeBool(buf, success)

	if success {
		writeInt32(buf, userID)
		writeString(buf, displayName)
		writeString(buf, statusText)

		if friends != nil {
			writeInt32(buf, len(friends))
			for _, f := range friends {
				writeInt32(buf, f.ID)
				writeString(buf, &f.Username)
				writeString(buf, f.DisplayName)
				writeInt32(buf, f.Status)
				writeString(buf, f.ImageURL)
				writeString(buf, f.StatusText)
			}
		} else {
			writeInt32(buf, 0)
		}

		if pendingContacts != nil {
			writeInt32(buf, len(pendingContacts))
			for _, c := range pendingContacts {
				writeInt32(buf, c.ID)
				writeString(buf, &c.Username)
				writeString(buf, c.DisplayName)
				writeString(buf, c.ImageURL)
				writeString(buf, c.Message)
			}
		} else {
			writeInt32(buf, len(pendingContacts))
		}
	}

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDLoginResult,
		Data: &bytes,
	}

	return &packet
}

// NewUserStatusChangePacket creates a new packet for user status change.
func NewUserStatusChangePacket(userID int, status int) *server.Packet {
	/*
		FriendUserID (int32)
		Status (int32)
	*/

	buf := new(bytes.Buffer)

	writeInt32(buf, userID)
	writeInt32(buf, status)

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDUserStatusChange,
		Data: &bytes,
	}

	return &packet
}

// NewChatFromPacket creates a new chat from user packet.
func NewChatFromPacket(fromUserID int, message string) *server.Packet {
	/*
		FromUserId (int32)
		MessageLen (int32)
		Message (string)
	*/
	buf := new(bytes.Buffer)

	writeInt32(buf, fromUserID)
	writeString(buf, &message)

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDChatFrom,
		Data: &bytes,
	}

	return &packet
}

// NewActionFromPacket creates a new action from user packet.
func NewActionFromPacket(fromUserID int, action string) *server.Packet {
	/*
		FromUserId (int32)
		ActionLen (int32)
		Action (string)
	*/
	buf := new(bytes.Buffer)

	writeInt32(buf, fromUserID)
	writeString(buf, &action)

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDActionFrom,
		Data: &bytes,
	}

	return &packet
}

// NewNudgeFromPacket creates a new nudge from user packet.
func NewNudgeFromPacket(fromUserID int) *server.Packet {
	/*
		FromUserId (int32)
	*/
	buf := new(bytes.Buffer)

	writeInt32(buf, fromUserID)

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDNudgeFrom,
		Data: &bytes,
	}

	return &packet
}

// NewAddContactResponsePacket creates a new add contact response packet.
func NewAddContactResponsePacket(resultCode int) *server.Packet {
	/*
		ResultCode (int32)
	*/
	buf := new(bytes.Buffer)

	writeInt32(buf, resultCode)

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDAddContactResponse,
		Data: &bytes,
	}

	return &packet
}

// NewNotifyAddRequestPacket creates a new notify add request packet.
func NewNotifyAddRequestPacket(userID int, username *string, displayName *string, message *string) *server.Packet {
	/*
		RequestedUserID (int32)
		UsernameLen (int32)
		Username (string)
		DisplayNameLen (int32)
		DisplayName (string)
		MessageLen (int32)
		Message (string)
	*/
	buf := new(bytes.Buffer)

	writeInt32(buf, userID)
	writeString(buf, username)
	writeString(buf, displayName)
	writeString(buf, message)

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDNotifyAddRequest,
		Data: &bytes,
	}

	return &packet
}

// NewConfirmContactResponsePacket creates a new confirm contact response packet.
func NewConfirmContactResponsePacket(resultCode int, requestedUserID int, user *models.UserModel) *server.Packet {
	/*
		ResultCode (jnt32)
		RequestedUserID (int32)

		(If successfull...)
		UsernameLen (int32)
		Username (string)
		DisplayNameLen (int32)
		DisplayName (string)
		Status (int32)
		ImageURLLen (int32)
		ImageURL (string)
		StatusTextLen (int32)
		StatusText (string)
	*/
	buf := new(bytes.Buffer)

	writeInt32(buf, resultCode)
	writeInt32(buf, requestedUserID)

	if user != nil {
		writeString(buf, &user.Username)
		writeString(buf, user.DisplayName)
		writeInt32(buf, user.Status)
		writeString(buf, user.ImageURL)
		writeString(buf, user.StatusText)
	}

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDConfirmContactResponse,
		Data: &bytes,
	}

	return &packet
}

// NewRejectContactResponsePacket creates a new reject contact response packet.
func NewRejectContactResponsePacket(resultCode int, requestedUserID int) *server.Packet {
	buf := new(bytes.Buffer)

	writeInt32(buf, resultCode)
	writeInt32(buf, requestedUserID)

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDRejectContactResponse,
		Data: &bytes,
	}

	return &packet
}

// NewAddContactAcceptedPacket creates a new contact accepted packet.
func NewAddContactAcceptedPacket(userID int, username *string, displayName *string, status int, imageURL *string, statusText *string) *server.Packet {
	buf := new(bytes.Buffer)

	writeInt32(buf, userID)
	writeString(buf, username)
	writeString(buf, displayName)
	writeInt32(buf, status)
	writeString(buf, imageURL)
	writeString(buf, statusText)

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDAddContactAccepted,
		Data: &bytes,
	}

	return &packet
}

// NewImageFromPacket creates a new image from user packet.
func NewImageFromPacket(fromUserID int, imageData *string) *server.Packet {
	/*
		FromUserId (int32)
		ImageDataLen (int32)
		ImageData (string)
	*/
	buf := new(bytes.Buffer)

	writeInt32(buf, fromUserID)
	writeString(buf, imageData)

	bytes := buf.Bytes()

	packet := server.Packet{
		ID:   PacketIDImageFrom,
		Data: &bytes,
	}

	return &packet
}
