package main

import (
	"flag"
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/junhsieh/goexamples/util"
	"github.com/junhsieh/iojson"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	sessionSecret     = "something-very-secret"
	sessionName       = "MySession"
	srvAddr           = "0.0.0.0:8443"
	srvWriteTimeout   = 15 * time.Second
	srvReadTimeout    = 15 * time.Second
	srvMaxHeaderBytes = 1 << 20
	CSRF_AUTH_KEY     = "31-byte-long-auth-key----------"
)

var (
	staticDir  = flag.String("staticDir", "/static/", "the directory to serve files from. Defaults to the static dir of the current dir")
	sessionDir = flag.String("sessionDir", "./session", "the directory to store sessions to. Defaults to the session dir of the current dir")
)

// TODO: include csrf for user authentication

// TODO: integrate with Redis. Store session in Redis.
// store session on server side.
var (
	store = sessions.NewFilesystemStore(*sessionDir, []byte(sessionSecret))
)

type (
	// MyHandlerV3 ...
	MyHandlerV3 func(http.ResponseWriter, *http.Request, *iojson.IOJSON)
)

func homeHTML() string {
	return `
		<html>
			<head>
				<script src="https://code.jquery.com/jquery-2.x-git.min.js"></script>
				<script src="static/debug.js?time=2"></script>
			</head>
			<body>
				<input type="text" id="Username" value="jun" placeholder="Username" />
				<br><input type="text" id="Password" value="junpass" placeholder="Password" />
				<br><select id="User">
					<option value="user1">User1</option>
					<option value="user2">User2</option>
					<option value="user3">User3</option>
				</select>
				<br><input type="text" id="CSRFToken" style="width: 800px;" />
				<br><button id="SalOrderBtn">SalOrder</button>
				<br><button id="NewsBtn">News</button>
				<br><button id="News3Btn">News 3</button>
				<br><button id="CSRFBtn">Get CSRF Token</button>
				<br><button id="LoginBtn">Login</button>
				<br><button id="LogoutBtn">Logout</button>
			</body>
		</html>
`
}

func srvHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("Host: %v", r.Host)
	// r.URL.Scheme will be empty if you're accessing the HTTP server not from an HTTP proxy,
	// a browser can issue a relative HTTP request instead of a absolute URL.
	// Additionally, you could check in the server/handler whether you get a
	// relative or absolute URL in the request by calling the IsAbs() method.
	// Reference: http://stackoverflow.com/questions/6899069/why-are-request-url-host-and-scheme-blank-in-the-development-server
	log.Printf("Scheme: %v", r.URL.Scheme)
	log.Printf("IsAbs: %v", r.URL.IsAbs())

	if _, err := w.Write([]byte(homeHTML())); err != nil {
		log.Printf("w.Write: %v", err.Error())
	}
}

func srvCSRFToken(w http.ResponseWriter, r *http.Request, o *iojson.IOJSON) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Get the token and pass it in the CSRF header. Our JSON-speaking client
	// or JavaScript framework can now read the header and return the token in
	// in its own "X-CSRF-Token" request header on the subsequent POST.
	w.Header().Set("X-CSRF-Token", csrf.Token(r))

}

func srvLogin(w http.ResponseWriter, r *http.Request, o *iojson.IOJSON) {
	// User input
	i := iojson.NewIOJSON()

	if err := i.Decode(r.Body); err != nil {
		o.AddError(err.Error())
		return
	}

	log.Printf("Username: %s", i.GetData("SQLUsername"))
	log.Printf("Password: %s", i.GetData("SQLPassword"))
	log.Printf("JSON: %s", i.EncodePretty())

	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, sessionName)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		o.AddError(err.Error())
		return
	}

	// TODO: add a logic to continue only when session.IsNew is true
	log.Printf("IsNew Session: %v", session.IsNew)

	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   "erp.local",
		MaxAge:   86400 * 1, // 3 Days
		Secure:   false,     // TODO: set to true once applied the SSL certificate.
		HttpOnly: true,
	}

	if i.GetData("SQLUsername") == "jun" && i.GetData("SQLPassword") == "junpass" {
		// Set some session values.
		session.Values["Username"] = i.GetData("SQLUsername")

		// Save it before we write to the response/return from the handler.
		if err := session.Save(r, w); err != nil {
			o.AddError(err.Error())
			return
		}
	}

	username := "init name"

	if s, ok := session.Values["Username"]; ok {
		username = s.(string)
	}

	o.AddData("test", "very good")
	o.AddData("username", username)
}

func srvLogout(w http.ResponseWriter, r *http.Request, o *iojson.IOJSON) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, sessionName)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		o.AddError(err.Error())
		return
	}

	// TODO: add a logic to continue only when session.IsNew is false
	log.Printf("IsNew Session: %v", session.IsNew)

	log.Printf("Domain: %v", session.Options.Domain)
	log.Printf("MaxAge Before: %v", session.Options.MaxAge)

	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   "erp.local",
		MaxAge:   -1,    // means delete cookie now.
		Secure:   false, // TODO: set to true once applied the SSL certificate.
		HttpOnly: true,
	}

	log.Printf("Domain: %v", session.Options.Domain)
	log.Printf("MaxAge After: %v", session.Options.MaxAge)

	// Save it before we write to the response/return from the handler.
	if err := session.Save(r, w); err != nil {
		o.AddError(err.Error())
		return
	}

	log.Printf("ID: %v", session.ID)

	if err := os.Remove("./session/session_" + session.ID); err != nil {
		//
	}

	o.AddData("msg", "cookie has been deleted from server")
}

func srvNotFound(w http.ResponseWriter, r *http.Request, o *iojson.IOJSON) {
	o.AddError("custom 404 page not found")
}

func getSession() {
}

func srvUserAuthentication(w http.ResponseWriter, r *http.Request, o *iojson.IOJSON) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, sessionName)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		o.AddError(err.Error())
		return
	}

	// TODO: how to determine whether a session is expired?

	// TODO: add a logic to continue only when session.IsNew is false (can be true too when the session is expired)
	log.Printf("IsNew Session: %v", session.IsNew)

	// TODO: need a way to check if session exists.

	if _, ok := session.Values["Username"]; !ok {
		//w.WriteHeader(http.StatusForbidden)
		o.AddError("You do not have the permission")
		return
	}

	o.AddData("welcome", "Welcom, "+session.Values["Username"].(string))

	// TODO: Add a logic to support both session/cookie and Header/Authorization
	/*
		Authorization := r.Header.Get("Authorization")

		log.Printf("Authorization: %v", Authorization)

		switch Authorization {
		case "user2":
			w.WriteHeader(http.StatusForbidden)
			o.AddError("You do not have the permission")
		case "user3":
			o.AddError("You do not have the permission")
		}
	*/
}

func srvSalOrder(w http.ResponseWriter, r *http.Request, o *iojson.IOJSON) {
	so := &struct {
		DealerName string
		IDShipAddr int64
		Price      float64
	}{}

	i := iojson.NewIOJSON()
	i.AddObj(so)

	if err := i.Decode(r.Body); err != nil {
		o.AddError(err.Error())
		return
	}

	log.Printf("SQLDealerName: %s", i.GetData("SQLDealerName"))
	log.Printf("SQLIDShipAddr: %d", int(i.GetData("SQLIDShipAddr").(float64)))
	log.Printf("SQLPrice: %f", i.GetData("SQLPrice"))
	log.Printf("JSON: %s", i.EncodePretty())

	vars := mux.Vars(r)
	IDSalOrder := vars["IDSalOrder"]

	o.AddData("IDSalOrder", IDSalOrder)
}

func srvNews(w http.ResponseWriter, r *http.Request, o *iojson.IOJSON) {
	news := struct {
		Subject string
		Author  string
		Body    string
	}{
		Subject: "Hello World",
		Author:  "Jun",
		Body:    "This is a Hello World message" + r.URL.Path,
	}

	o.AddObj(news)

}

// version 1: returning http.HandlerFunc
func srvNoCheckingHandlerFunc(mh ...MyHandlerV3) http.HandlerFunc {
	return handlerLoopHandlerFunc(append([]MyHandlerV3{}, mh...))
}

// version 1: returning http.HandlerFunc
func srvRegularCheckingHandlerFunc(mh ...MyHandlerV3) http.HandlerFunc {
	return handlerLoopHandlerFunc(append([]MyHandlerV3{srvUserAuthentication}, mh...))
}

// version 1: returning http.HandlerFunc
func handlerLoopHandlerFunc(myhandlers []MyHandlerV3) http.HandlerFunc {
	// TODO: benchmark returning http.Handler vs http.HandlerFunc
	return gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := iojson.NewIOJSON()

		for _, myhandler := range myhandlers {
			if myhandler(w, r, o); o.ErrCount > 0 {
				o.Echo(w)
				return
			}
		}

		o.Echo(w)
	})).(http.HandlerFunc)
}

// version 2: returning http.Handler
func srvNoCheckingHandler(mh ...MyHandlerV3) http.Handler {
	//return srvCSRF(handlerLoopHandler(append([]MyHandlerV3{}, mh...)))
	return srvGZIPCSRF(handlerLoopHandler(append([]MyHandlerV3{}, mh...)))
}

func srvGZIPCSRF(h http.Handler) http.Handler {
	return srvGZIP(srvCSRF(h))
}

func srvCSRF(h http.Handler) http.Handler {
	return csrf.Protect([]byte(CSRF_AUTH_KEY),
		csrf.Secure(true),
		csrf.MaxAge(3600*12),
		csrf.ErrorHandler(http.HandlerFunc(unauthorizedHandler)),
	)(h)
}

func srvGZIP(h http.Handler) http.Handler {
	return gziphandler.GzipHandler(h)
}

// version 2: returning http.Handler
func handlerLoopHandler(myhandlers []MyHandlerV3) http.Handler {
	// TODO: benchmark returning http.Handler vs http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := iojson.NewIOJSON()

		for _, myhandler := range myhandlers {
			if myhandler(w, r, o); o.ErrCount > 0 {
				o.Echo(w)
				return
			}
		}

		o.Echo(w)
	})
}

//
func unauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	o := iojson.NewIOJSON()
	o.AddError("Forbidden - CSRF token invalid")

	o.Echo(w)
}

func main() {
	//

	flag.Parse()

	//
	util.HelpGenTLSKeys()

	//
	r := mux.NewRouter()
	r.Host("erp.local:8080")

	// custom 404 not found handler
	//r.NotFoundHandler = http.HandlerFunc(srvNotFound) // NOTE: standware way.
	r.NotFoundHandler = srvNoCheckingHandlerFunc(srvNotFound)

	//
	r.HandleFunc("/", srvHome)

	r.HandleFunc("/Login", srvNoCheckingHandlerFunc(srvLogin))
	r.HandleFunc("/Logout", srvNoCheckingHandlerFunc(srvLogout))

	r.HandleFunc("/News1", srvNoCheckingHandlerFunc(srvNews))
	r.Handle("/News2", srvNoCheckingHandler(srvNews))

	r.HandleFunc("/SalOrder/{IDSalOrder}", srvRegularCheckingHandlerFunc(srvSalOrder))

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix(*staticDir).Handler(http.StripPrefix(*staticDir, http.FileServer(http.Dir("."+*staticDir))))

	// Subrouter. Only matches if domain is "www.example.com".
	//s := r.Host("erp.local").Subrouter()
	//s.HandleFunc("/SalOrder/{IDSalOrder}", srvRegularCheckingHandlerFunc(srvSalOrder))

	// Subrouter.
	s2 := r.Host("erp.local").Subrouter()
	s2.Handle("/CSRFToken", srvNoCheckingHandler(srvCSRFToken))
	s2.Handle("/News3", srvNoCheckingHandler(srvNews))

	//
	srv := &http.Server{
		Handler: r,
		Addr:    srvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:   srvWriteTimeout,
		ReadTimeout:    srvReadTimeout,
		MaxHeaderBytes: srvMaxHeaderBytes,
	}

	//if err := srv.ListenAndServe(); err != nil {
	if err := srv.ListenAndServeTLS("mydomain.com.crt", "mydomain.com.key"); err != nil {
		log.Printf("srv.ListenAndServe: %v", err.Error())
	}
}
