package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
)

var (
	certFile = flag.String("cert", "client-crt.pem", "client certificate")
	keyFile  = flag.String("key", "client-privetekey.pem", "client private key")
	caFile   = flag.String("ca", "ca-cert.pem", "pem eoncoded CA certificate")
	port     = flag.String("port", "5001", "server port")
	host     = flag.String("host", "127.0.0.1", "server host")
)

func main() {
	flag.Parse()

	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cfg := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	cfg.BuildNameToCertificate()

	conn, err := tls.Dial("tcp", *host+":"+*port, cfg)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 100)
	if _, err := conn.Read(buf); err != nil {
		log.Fatal(err)
	}
	if _, err := conn.Write([]byte("hey\n")); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", buf)
}
