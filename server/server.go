// Package server provides a multi-client TCP socket server. It doesn't know anything beyond a packet
// of data that can be sent/received in TLV format (type-length-value).
package server

import (
	"log"
	"net"
	"sync"
	"time"
)

// OnHandlePacket is used as an event for when a packet is received from a client.
type OnHandlePacket func(s *TCPServer, c net.Conn, p *Packet)

// OnClientConnect is used as event for when a new client has connected.
type OnClientConnect func(s *TCPServer, c net.Conn, addr string)

// OnClientDisconnect is used as event for when a client has disconnected.
type OnClientDisconnect func(s *TCPServer, c net.Conn, addr string)

// TCPServer represents the listening socket and all connected clients.
type TCPServer struct {
	listener           net.Listener
	clients            []*client
	mutex              *sync.Mutex
	onHandlePacket     OnHandlePacket
	onClientConnect    OnClientConnect
	onClientDisconnect OnClientDisconnect
}

// NewTCPServer creates a new TCPServer instance.
func NewTCPServer(h OnHandlePacket, c OnClientConnect, d OnClientDisconnect) *TCPServer {
	return &TCPServer{
		mutex:              &sync.Mutex{},
		onHandlePacket:     h,
		onClientConnect:    c,
		onClientDisconnect: d,
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
			bytes := p.toBytes()

			client.conn.SetWriteDeadline(time.Now().Add(writeTimeoutDuration))
			client.conn.Write(*bytes)
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

	s.onClientDisconnect(s, client.conn, client.conn.RemoteAddr().String())

	client.conn.Close()
}

func (s *TCPServer) serve(client *client) {
	defer s.remove(client)

	for {
		packet, err := client.readPacket()
		if err != nil {
			break
		}

		s.onHandlePacket(s, client.conn, packet)
	}
}
