package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

func createClientConfig(ca, crt, key string) (*tls.Config, error) {
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
		RootCAs:      roots,
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
	log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")
}

func main() {
	connect := flag.String("connect", "example.com:4433", "who to connect to")
	ca := flag.String("ca", "./ca.crt", "root certificate")
	crt := flag.String("crt", "./client.crt", "certificate")
	key := flag.String("key", "./client.key", "key")
	flag.Parse()

	addr := *connect
	if !strings.Contains(addr, ":") {
		addr += ":443"
	}

	config, err := createClientConfig(*ca, *crt, *key)
	if err != nil {
		log.Fatal("config failed: %s", err.Error())
	}

	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("failed to connect: %s", err.Error())
	}
	defer conn.Close()

	log.Printf("connect to %s succeed", addr)
	printConnState(conn)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		io.Copy(conn, os.Stdin)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		io.Copy(os.Stdout, conn)
		wg.Done()
	}()
	wg.Wait()
}
