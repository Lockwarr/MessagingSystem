package server

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"

	protocol "../protocol"
)

type client struct {
	conn      net.Conn
	publicKey string
	writer    *protocol.CommandWriter
}

// TCPChatServer - struct for the tcp server
type TCPChatServer struct {
	listener net.Listener
	clients  []*client
	mutex    *sync.Mutex
}

// UnknownClient - custom error for unknown clients
var (
	UnknownClient = errors.New("Unknown client")
)

// NewServer - used for creating new TCP chat server
func NewServer() *TCPChatServer {
	return &TCPChatServer{
		mutex: &sync.Mutex{},
	}
}

// Listen - listens on given adrdress
func (s *TCPChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)

	if err == nil {
		s.listener = l
	}

	log.Printf("Listening on %v", address)

	return err
}

// Close - closes sever listener
func (s *TCPChatServer) Close() {
	s.listener.Close()
}

// Start - starts the server and accept connections
func (s *TCPChatServer) Start() {
	for {
		// XXX: need a way to break the loop
		conn, err := s.listener.Accept()

		if err != nil {
			log.Print(err)
		} else {
			// handle connection
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

// Broadcast - broadcasts to other server clients
func (s *TCPChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		// TODO: handle error here?
		client.writer.Write(command)
	}

	return nil
}

// Send - sends message to other clients
func (s *TCPChatServer) Send(publicKey string, command interface{}) error {
	for _, client := range s.clients {
		if client.publicKey == publicKey {
			return client.writer.Write(command)
		}
	}

	return UnknownClient
}

func (s *TCPChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}

	s.clients = append(s.clients, client)

	return client
}

func (s *TCPChatServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// remove the connections from clients array
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}

	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

func (s *TCPChatServer) serve(client *client) {
	cmdReader := protocol.NewCommandReader(client.conn)

	defer s.remove(client)

	for {
		cmd, err := cmdReader.Read()

		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.SendCommand:
				go s.Broadcast(protocol.MessageCommand{
					Message:   v.Message,
					PublicKey: client.publicKey,
				})
			}
		}

		if err == io.EOF {
			break
		}
	}
}
