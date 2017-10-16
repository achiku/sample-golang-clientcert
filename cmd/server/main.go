package main

import (
	"crypto/tls"
	"flag"
	"log"
)

var (
	serverCrt = flag.String("cert", "server.crt", "server certificate")
	serverKey = flag.String("key", "server.key", "server private key")
	port      = flag.String("port", "5001", "server port")
)

func main() {
	flag.Parse()

	cer, err := tls.LoadX509KeyPair(*serverCrt, *serverKey)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cer},
	}
	addr := "127.0.0.1:" + *port
	s := NewTCPServer(addr)

	s.OnNewClient(func(c *Client) {
		c.Send("hello from TCP server")
	})
	s.OnNewMessage(func(c *Client, msg string) {
		log.Printf("new message: %s", msg)
	})
	s.OnClientConnectionClosed(func(c *Client, err error) {
		log.Print("client connection closed")
	})

	log.Printf("server start at %s", addr)
	if err := s.Listen(cfg); err != nil {
		log.Fatal(err)
	}
}
