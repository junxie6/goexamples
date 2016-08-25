package main

import (
	"log"
	"net/http"
	"time"
)

const (
	srvAddr           = "0.0.0.0:8443"
	srvWriteTimeout   = 15 * time.Second
	srvReadTimeout    = 15 * time.Second
	srvMaxHeaderBytes = 1 << 20
)

func srvHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL %v", r.URL.Path)

	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	// Reference:
	// https://golang.org/pkg/net/http/#example_ServeMux_Handle
	// https://golang.org/pkg/net/http/#ServeMux
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello World!"))
}

func main() {
	HelpGenTLSKeys()

	// Reference:
	// https://golang.org/pkg/net/http/#ServeMux
	// https://golang.org/pkg/net/http/#example_ServeMux_Handle
	// http://golang.org/pkg/net/http/#Server
	mux := http.NewServeMux()

	mux.HandleFunc("/", srvHome)

	//
	srv := &http.Server{
		Handler: mux,
		Addr:    srvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:   srvWriteTimeout,
		ReadTimeout:    srvReadTimeout,
		MaxHeaderBytes: srvMaxHeaderBytes,
	}

	log.Printf("Listening to TCP address: %v", srv.Addr)

	//
	//if err := srv.ListenAndServe(); err != nil {
	if err := srv.ListenAndServeTLS("mydomain.com.crt", "mydomain.com.key"); err != nil {
		log.Printf("srv.ListenAndServe: %v", err.Error())
	}
}

func HelpGenTLSKeys() {
	str := `
To generate the private key and the self-signed certificate:

Use this method if you want to use HTTPS (HTTP over TLS) to secure your Apache HTTP or Nginx web server, and you want to use a Certificate Authority (CA) to issue the SSL certificate. The CSR that is generated can be sent to a CA to request the issuance of a CA-signed SSL certificate. If your CA supports SHA-2, add the -sha256 option to sign the CSR with SHA-2.

# openssl req -newkey rsa:2048 -nodes -subj "/C=CA/ST=British Columbia/L=Vancouver/O=My Company Name/CN=mydomain.com" -keyout mydomain.com.key -out mydomain.com.csr

Note: The -newkey rsa:2048 option specifies that the key should be 2048-bit, generated using the RSA algorithm.
Note: The -nodes option specifies that the private key should not be encrypted with a pass phrase.
Note: The -new option, which is not included here but implied, indicates that a CSR is being generated.

Generate a Self-Signed Certificate:

Use this method if you want to use HTTPS (HTTP over TLS) to secure your Apache HTTP or Nginx web server, and you do not require that your certificate is signed by a CA.

This command creates a 2048-bit private key (domain.key) and a self-signed certificate (domain.crt) from scratch:

# openssl req -newkey rsa:2048 -nodes -subj "/C=CA/ST=British Columbia/L=Vancouver/O=My Company Name/CN=mydomain.com" -keyout mydomain.com.key -x509 -days 365 -out mydomain.com.crt
`
	log.Printf("%v\n", str)
}
