package main

import (
	"flag"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/junhsieh/alice"
	"github.com/junhsieh/goexamples/util"
	"github.com/junhsieh/gziphandler"
	"github.com/junhsieh/iojson"
	"github.com/junhsieh/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	srvWriteTimeout   = 15 * time.Second
	srvReadTimeout    = 15 * time.Second
	srvMaxHeaderBytes = 1 << 20
)

var (
	srvDomain     = flag.String("srvDomain", "erp.local", `Server's domain name`)
	srvPort       = flag.String("srvPort", ":8443", `Server's port number`)
	certTLS       = flag.String("certTLS", "mydomain.com.crt", `TLS certificate`)
	keyTLS        = flag.String("keyTLS", "mydomain.com.key", `TLS key`)
	staticDir     = flag.String("staticDir", "/static/", `The directory to serve files from. Defaults to the static dir of the current dir`)
	sessionDir    = flag.String("sessionDir", "./session", `The directory to store sessions to. Defaults to the session dir of the current dir`)
	sessionName   = flag.String("sessionName", "MySession", `Session name`)
	sessionSecret = flag.String("sessionSecret", "something-very-secret", `Session secret`)
	authCSRFKey   = flag.String("authCSRFKey", "31-byte-long-auth-key----------", `CSRF Auth key`)
)

// TODO: include csrf for user authentication

// TODO: integrate with Redis. Store session in Redis.
// store session on server side.
var (
	store = sessions.NewFilesystemStore(*sessionDir, []byte(*sessionSecret))
)

type (
	// MyHandlerV3 ...
	MyHandlerV3 func(http.ResponseWriter, *http.Request, *iojson.IOJSON)
)

func homeHTML() string {
	return `
		<html>
			<head>
				<link rel="shortcut icon" href="data:image/x-icon;," type="image/x-icon">
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
	session, err := store.Get(r, *sessionName)

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
		Domain:   *srvDomain,
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
	session, err := store.Get(r, *sessionName)

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
		Domain:   *srvDomain,
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

func srvUserAuthentication(w http.ResponseWriter, r *http.Request, o *iojson.IOJSON) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, *sessionName)

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

//
func unauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	o := iojson.NewIOJSON()
	o.AddError("Forbidden - CSRF token invalid")

	o.Echo(w)
}

func srvHome2(w http.ResponseWriter, r *http.Request) {
	//log.Printf("Inside home2")

	o := r.Context().Value("iojson").(*iojson.IOJSON)
	o.AddData("act", "home 2; "+r.FormValue("act"))
}

func main() {
	//
	flag.Parse()

	//
	util.HelpGenTLSKeys()

	//
	chain := alice.New(
		gziphandler.GzipHandler,
		iojson.EchoHandler,
		middleware.LoggerHandler,
		middleware.DomainHandler(*srvDomain+*srvPort),
	)

	//
	mux := http.NewServeMux()

	// NOTE: FileServer calls path.Clean() to clean up path.
	mux.Handle(*staticDir, http.StripPrefix(*staticDir, http.FileServer(http.Dir("."+*staticDir))))

	mux.Handle("/Home2", chain.Then(http.HandlerFunc(srvHome2)))

	//
	srv := &http.Server{
		Handler: mux,
		Addr:    *srvPort,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout:   srvWriteTimeout,
		ReadTimeout:    srvReadTimeout,
		MaxHeaderBytes: srvMaxHeaderBytes,
	}

	//if err := srv.ListenAndServe(); err != nil {
	if err := srv.ListenAndServeTLS(*certTLS, *keyTLS); err != nil {
		log.Printf("srv.ListenAndServe: %v", err.Error())
	}
}
