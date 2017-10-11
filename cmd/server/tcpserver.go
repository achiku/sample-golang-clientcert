package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
)

// Client holds info about connection
type Client struct {
	conn   net.Conn
	Server *Server
}

// Server TCP server
type Server struct {
	address                  string // Address to open connection: localhost:9999
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

// Read client data from channel
func (c *Client) listen() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, message)
	}
}

// Send text message to client
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// SendBytes bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

// Conn connection
func (c *Client) Conn() net.Conn {
	return c.conn
}

// Close close connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// OnNewClient Called right after server starts listening new client
func (s *Server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// OnClientConnectionClosed Called right after connection closed
func (s *Server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// OnNewMessage Called when Client receives new message
func (s *Server) OnNewMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

// Listen network server
func (s *Server) Listen(tlsCfg *tls.Config) error {
	var (
		listener net.Listener
		err      error
	)
	if tlsCfg == nil {
		listener, err = net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
	} else {
		listener, err = tls.Listen("tcp", s.address, tlsCfg)
		if err != nil {
			return err
		}
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.listen()
		s.onNewClientCallback(client)
	}
}

// NewTCPServer creates new tcp server instance
func NewTCPServer(address string) *Server {
	server := &Server{
		address: address,
	}
	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})
	return server
}
