package client

import (
	"io"
	"log"
	"net"

	helpers "../helpers"
	protocol "../protocol"
)

// TCPChatClient - struct for the chat clientee
type TCPChatClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	publicKey string
	incoming  chan protocol.MessageCommand
}

// NewClient - creating new client and generating public/private key
func NewClient() *TCPChatClient {
	return &TCPChatClient{
		incoming:  make(chan protocol.MessageCommand),
		publicKey: helpers.GeneratePublicKey(),
	}
}

// Dial - Dials to given address
func (c *TCPChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)

	if err == nil {
		c.conn = conn
	}

	c.cmdReader = protocol.NewCommandReader(conn)
	c.cmdWriter = protocol.NewCommandWriter(conn)

	return err
}

// Start - reads incoming messages from the server or other clients
func (c *TCPChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Read error %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			default:
				log.Printf("Unknown command: %v", v)
			}
		}
	}
}

// Close - closes the connection of the client
func (c *TCPChatClient) Close() {
	c.conn.Close()
}

// SetName - sets the name of the user
func (c *TCPChatClient) SetName(name string) error {
	return c.Send(protocol.NameCommand{
		Name: name,
	})
}

// Incoming - returns incoming messages
func (c *TCPChatClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

// Send - sends command
func (c *TCPChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

// SendMessage - sends a message command
func (c *TCPChatClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
	})
}
