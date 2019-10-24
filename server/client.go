package server

import (
	"io"
	"log"
	"net"
	"time"
)

// A wrapper for connected clients as later we will store more information around a connected socket.
type client struct {
	conn net.Conn
	name string
}

func (c *client) readPacket() (packet *Packet, err error) {

	packetType, e := readFourBytes(c)
	if e != nil {
		return nil, e
	}

	packetLength, e := readFourBytes(c)
	if e != nil {
		return nil, e
	}

	packetData, e := readPacketData(c, packetLength)
	if e != nil {
		return nil, e
	}

	return &Packet{
		ID:   packetType,
		Data: packetData,
	}, nil
}

func readFourBytes(c *client) (val int64, err error) {
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

func readPacketData(c *client, length int64) (data *[]byte, err error) {
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
