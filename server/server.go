// Package server provides a multi-client TCP socket server. It doesn't know anything beyond a packet
// of data that can be sent/received in TLV format (type-length-value).
package server

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Packet represents a received packet from a client connection.
type Packet struct {
	ID   int64
	Data []byte
}

// OnHandlePacket is used as an event for when a packet is received from a client.
type OnHandlePacket func(s *TCPServer, c net.Conn, p Packet)

// OnClientConnect is used as event for when a new client has connected.
type OnClientConnect func(s *TCPServer, c net.Conn, addr string)

// A wrapper for connected clients as later we will store more information around a connected socket.
type client struct {
	conn net.Conn
	name string
}

// TCPServer represents the listening socket and all connected clients.
type TCPServer struct {
	listener        net.Listener
	clients         []*client
	mutex           *sync.Mutex
	onHandlePacket  OnHandlePacket
	onClientConnect OnClientConnect
}

// NewTCPServer creates a new TCPServer instance.
func NewTCPServer(h OnHandlePacket, c OnClientConnect) *TCPServer {
	return &TCPServer{
		mutex:           &sync.Mutex{},
		onHandlePacket:  h,
		onClientConnect: c,
	}
}

// Listen opens a port for connections on the specified address.
func (s *TCPServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)

	if err == nil {
		s.listener = l
	}

	log.Printf("Listening on %v", address)

	return err
}

// Close the listening server socket.
func (s *TCPServer) Close() {
	s.listener.Close()
}

// Run start an indefinite loop waiting for new clients to connect and serves them. Messages from clients
// are sent through the TCPServer packet channel.
func (s *TCPServer) Run() {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Print(err)
		} else {
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

// BroadcastPacket will send the specified packet to all connected clients except the specified client connection.
func (s *TCPServer) BroadcastPacket(p Packet, c net.Conn) {
	for _, client := range s.clients {
		if client.conn != c {
			packetIDBytes := make([]byte, 4)
			packetLengthBytes := make([]byte, 4)

			binary.LittleEndian.PutUint32(packetIDBytes, uint32(p.ID))
			binary.LittleEndian.PutUint32(packetLengthBytes, uint32(len(p.Data)))

			buffer := append(packetIDBytes, packetLengthBytes...)
			buffer = append(buffer, p.Data...)

			timeoutDuration := 30 * time.Second
			client.conn.SetWriteDeadline(time.Now().Add(timeoutDuration))

			client.conn.Write(buffer)
		}
	}
}

func (s *TCPServer) accept(c net.Conn) *client {
	s.onClientConnect(s, c, c.RemoteAddr().String())

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn: c,
	}

	s.clients = append(s.clients, client)

	return client
}

func (s *TCPServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}

	log.Printf(
		"Closing connection from %v",
		client.conn.RemoteAddr().String())

	// TODO: Raise a client disconnected event here.

	client.conn.Close()
}

func (s *TCPServer) serve(client *client) {
	defer s.remove(client)

	timeoutDuration := 60 * time.Second

	for {
		// TODO: Move this in to some kind of packet parser package...

		packetTypeBytes := make([]byte, 4)
		packetLengthBytes := make([]byte, 4)

		client.conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		intReader := io.LimitReader(client.conn, 4)
		_, err := intReader.Read(packetTypeBytes)
		if err != nil {
			log.Printf("Packet type read error: %v", err)
			break
		}

		client.conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		intReader = io.LimitReader(client.conn, 4)
		_, err = intReader.Read(packetLengthBytes)
		if err != nil {
			log.Printf("Packet length read error: %v", err)
			break
		}

		// Packet header values are 32-bit, so convert them.
		packetType := bytesToInt64(packetTypeBytes)
		packetLength := bytesToInt64(packetLengthBytes)

		packetData := make([]byte, packetLength)
		valueReader := io.LimitReader(client.conn, packetLength)

		// Read the packet body.
		if packetLength > 0 {
			client.conn.SetReadDeadline(time.Now().Add(timeoutDuration))

			_, err = valueReader.Read(packetData)
			if err != nil {
				log.Printf("Packet data read error: %v", err)
				break
			}
		}

		s.onHandlePacket(s, client.conn, Packet{
			ID:   packetType,
			Data: packetData,
		})
	}
}

func bytesToInt64(b []byte) int64 {
	var val int64

	val |= int64(b[0])
	val |= int64(b[1]) << 8
	val |= int64(b[2]) << 16
	val |= int64(b[3]) << 24

	return val
}
