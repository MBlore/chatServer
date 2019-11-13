package server

import (
	"errors"
	"io"
	"log"
	"net"
	"time"
)

// Client represents a connected user.
type Client struct {
	conn        net.Conn
	Username    string
	DisplayName *string
	LoggedIn    bool
	UserID      int
}

// SendPacket sends the specified packet to the connected client.
func (c *Client) SendPacket(p *Packet) error {
	bytes := p.toBytes()

	c.conn.SetWriteDeadline(time.Now().Add(writeTimeoutDuration))
	_, err := c.conn.Write(*bytes)

	return err
}

func (c *Client) readPacket() (packet *Packet, err error) {

	packetType, e := readFourBytes(c)
	if e != nil {
		return nil, e
	}

	if packetType < 0 {
		return nil, errors.New("invalid packet id " + string(packetType))
	}

	packetLength, e := readFourBytes(c)
	if e != nil {
		return nil, e
	}

	if packetLength > maxPacketLength {
		return nil, errors.New("packet length above maximum allowed")
	}

	var packetData *[]byte = nil

	if packetLength > 0 {
		data, e := readPacketData(c, packetLength)
		if e != nil {
			return nil, e
		}

		packetData = data
	}

	return &Packet{
		ID:   packetType,
		Data: packetData,
	}, nil
}

func readFourBytes(c *Client) (val int64, err error) {
	bytes := make([]byte, 4)

	c.conn.SetReadDeadline(time.Now().Add(readTimeoutDuration))

	reader := io.LimitReader(c.conn, 4)
	_, e := reader.Read(bytes)

	if e != nil {
		log.Printf("Read error: %v", e)
		return 0, e
	}

	return bytesToInt64(bytes), nil
}

func readPacketData(c *Client, length int64) (data *[]byte, err error) {
	if length == 0 {
		return nil, nil
	}

	bytes := make([]byte, length)

	c.conn.SetReadDeadline(time.Now().Add(readTimeoutDuration))

	reader := io.LimitReader(c.conn, length)

	_, e := reader.Read(bytes)

	if e != nil {
		log.Printf("Read error: %v", e)
		return nil, e
	}

	return &bytes, nil
}
