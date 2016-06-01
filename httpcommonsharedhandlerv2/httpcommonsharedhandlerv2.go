package main

import (
	"errors"
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
	"strings"
)

// Reference:
// http://stackoverflow.com/questions/26204485/gorilla-mux-custom-middleware
// http://laicos.com/writing-handsome-golang-middleware/
// https://justinas.org/alice-painless-middleware-chaining-for-go/
// https://www.nicolasmerouze.com/middlewares-golang-best-practices-examples/
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81#.4xnxwa6mx

type (
	MyHandler func(w http.ResponseWriter, r *http.Request) error
)

func serveHandlers(myhandlers ...MyHandler) http.Handler {
	return gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, myhandler := range myhandlers {
			if err := myhandler(w, r); err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		}
	}))
}

func serveLog(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Log IP: %s\n", r.Header.Get("X-Forwarded-For"))
	return nil
}

func serveCheckAllowHost(w http.ResponseWriter, r *http.Request) error {
	log.Printf("Check IP: %s\n", r.Header.Get("X-Forwarded-For"))

	if ok := isAllowedHost(r.RemoteAddr, r.Header.Get("X-Forwarded-For")); !ok {
		return errors.New("Disallow IP")
	}

	return nil
}

func serveTest1(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Test1"))
	return nil
}

func serveTest2(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Test2"))
	return nil
}

func serveTest3(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("Test3"))
	return nil
}

func isAllowedHost(proxyIPPort string, userIP string) bool {
	// In case of IPv4 most client addresses are masked behind NAT, on your server side you ONLY see the globally routable address
	// which is the router's own global address.
	// Be caution against using the X-Forwarded-For header for ANYTHING unless it comes from a trusted source (e.g. your own reverse proxy).
	// The client can set this header to an arbitrary value and can cause some funny or even dangerous bugs to be triggered.
	allowProxyIP := "192.168.0.1"

	//test := proxy_ip[0:strings.Index(proxy_ip, ":")]
	proxyIPPortArr := strings.Split(proxyIPPort, ":")

	if proxyIPPortArr[0] != allowProxyIP {
		log.Printf("Disallow Proxy IP: %s", proxyIPPortArr[0])
		return false
	}

	//return true // remove the user's IP restriction

	allowHostArr := []string{
		"192.168.0.2",
	}

	for _, v := range allowHostArr {
		if userIP == v {
			return true
		}
	}

	log.Printf("Disallow user IP: %s", userIP)
	return false
}

func main() {
	commonHandlers := []MyHandler{serveLog, serveCheckAllowHost}

	http.Handle("/test", serveHandlers(append(commonHandlers, serveTest1, serveTest2)...))
	http.Handle("/test3", serveHandlers(append(commonHandlers, serveTest3)...))

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
