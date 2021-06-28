package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
)

func createServerConfig(ca, crt, key string) (*tls.Config, error) {
	caCertPEM, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, err
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
	}, nil
}

func printConnState(conn *tls.Conn) {
	log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")
	state := conn.ConnectionState()
	log.Printf("Version: %x", state.Version)
	log.Printf("HandshakeComplete: %t", state.HandshakeComplete)
	log.Printf("DidResume: %t", state.DidResume)
	log.Printf("CipherSuite: %x", state.CipherSuite)
	log.Printf("NegotiatedProtocol: %s", state.NegotiatedProtocol)
	log.Printf("NegotiatedProtocolIsMutual: %t", state.NegotiatedProtocolIsMutual)

	log.Print("Certificate chain:")
	for i, cert := range state.PeerCertificates {
		subject := cert.Subject
		issuer := cert.Issuer
		log.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
		log.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	}
	//log.Printf("%#v", state)
	log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")
}

func main() {
	listen := flag.String("listen", "localhost:4433", "which port to listen")
	ca := flag.String("ca", "./ca.crt", "root certificate")
	crt := flag.String("crt", "./server.crt", "certificate")
	key := flag.String("key", "./server.key", "key")
	flag.Parse()

	config, err := createServerConfig(*ca, *crt, *key)
	if err != nil {
		log.Fatal("config failed: %s", err.Error())
	}

	ln, err := tls.Listen("tcp", *listen, config)
	if err != nil {
		log.Fatal("listen failed: %s", err.Error())
	}

	log.Printf("listen on %s", *listen)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("accept failed: %s", err.Error())
			break
		}
		log.Printf("connection open: %s", conn.RemoteAddr())
		//printConnState(conn.(*tls.Conn))

		go func(c net.Conn) {
			wr, _ := io.Copy(c, c)
			printConnState(conn.(*tls.Conn))
			c.Close()
			log.Printf("connection close: %s, written: %d", conn.RemoteAddr(), wr)
		}(conn)
	}
}
