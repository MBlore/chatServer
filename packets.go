package main

import (
	"bytes"
	"chatServer/server"
	"encoding/binary"
)

// Packet types supported in our client/server protocol.
const (
	PacketIDLogin            = 0
	PacketIDLoginResult      = 1
	PacketIDChat             = 2
	PacketIDNudge            = 3
	PacketIDAudio            = 4
	PacketIDPing             = 5
	PacketIDAction           = 6
	PacketIDActionFrom       = 7
	PacketIDSetDisplayName   = 8
	PacketIDNudgeFrom        = 9
	PacketIDUserStatusChange = 10
	PacketIDChatFrom         = 11
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
func NewLoginResultPacket(success bool, userID int, displayName *string, friends []FriendModel) *server.Packet {
	/*
		Success (byte)
		UserId (int32)
		DisplayNameLength (int32)
		DisplayName (string)
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
	*/

	buf := new(bytes.Buffer)

	writeBool(buf, success)

	if success {
		writeInt32(buf, userID)
		writeString(buf, displayName)

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
