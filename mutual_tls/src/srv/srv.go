package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Edit hosts file:
// # vim /etc/hosts
//
// 127.0.0.1       example.com
//
// Generate CA private key:
//
// $ mkdir key
// $ openssl genrsa -out key/ca.key 2048
//
// Genearte CA certificate:
//
// $ openssl req -new -x509 -days 365 -key key/ca.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out key/ca.crt
//
// Generate server's private key and certificate sign request:
//
// $ openssl req -newkey rsa:2048 -nodes -keyout key/server.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=*.example.com" -out key/server.csr
//
// Generate server's certificate:
//
// $ openssl x509 -req -extfile <(printf "subjectAltName=DNS:example.com,DNS:www.example.com") -days 365 -in key/server.csr -CA key/ca.crt -CAkey key/ca.key -CAcreateserial -out key/server.crt
//
// Generate client's private key and certificate sign request:
//
// $ openssl req -newkey rsa:2048 -nodes -keyout key/client.key -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=client.example.com" -out key/client.csr
//
// Generate client's certificate:
//
// $ openssl x509 -req -extfile <(printf "subjectAltName=DNS:client.example.com") -days 365 -in key/client.csr -CA key/ca.crt -CAkey key/ca.key -CAcreateserial -out key/client.crt
//
// Check X509v3 Subject Alternative Name
//
// $ openssl x509 -in server.crt -text -noout
// or
// $ openssl req -in server.csr -text -noout
//
// Reference:
// https://golang.org/pkg/crypto/x509/
// https://golang.org/pkg/crypto/tls/
// https://venilnoronha.io/a-step-by-step-guide-to-mtls-in-go
// https://smallstep.com/hello-mtls/doc/server/go
// https://smallstep.com/hello-mtls/doc/client/go
// https://github.com/smallstep/cli#installation-guide

func helloHandler(w http.ResponseWriter, r *http.Request) {
	printConnState(r.TLS)

	// Write "Hello, world!" to the response body
	io.WriteString(w, "Hello, world!\n")
}

func printConnState(state *tls.ConnectionState) {
	log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")
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
		log.Printf(" %d subject:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
		log.Printf("   issuer:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	}
	//log.Printf("%#v", state)
	log.Print(">>>>>>>>>>>>>>>> State End <<<<<<<<<<<<<<<<")
}

func main() {
	// Set up a /hello resource handler
	http.HandleFunc("/hello", helloHandler)

	// Create a CA certificate pool and add ca.crt to it
	caCert, err := ioutil.ReadFile("key/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	//tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS("key/server.crt", "key/server.key"))
}
