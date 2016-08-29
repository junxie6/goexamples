package main

import (
	"flag"
	"github.com/gorilla/csrf"
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
	srvDomain       = flag.String("srvDomain", "erp.local", `Server's domain name`)
	srvPort         = flag.String("srvPort", ":8443", `Server's port number`)
	certTLS         = flag.String("certTLS", "erp.local.crt", `TLS certificate`)
	keyTLS          = flag.String("keyTLS", "erp.local.key", `TLS key`)
	staticDir       = flag.String("staticDir", "/static/", `The directory to serve files from. Defaults to the static dir of the current dir`)
	sessionDir      = flag.String("sessionDir", "./session", `The directory to store sessions to. Defaults to the session dir of the current dir`)
	sessionName     = flag.String("sessionName", "MySession", `Session name`)
	sessionSecret   = flag.String("sessionSecret", "something-very-secret", `Session secret`)
	sessionMaxAge   = flag.Int("sessionMaxAge", 3600*12, `Session MaxAge`)
	sessionSecure   = flag.Bool("sessionSecure", true, `Session Secure`)
	sessionHttpOnly = flag.Bool("sessionHttpOnly", true, `Session HttpOnly`)
	csrfAuthKey     = flag.String("csrfAuthKey", "31-byte-long-auth-key----------", `CSRF Auth key`)
	csrfSecure      = flag.Bool("csrfSecure", true, `CSRF Secure`)
	csrfMaxAge      = flag.Int("csrfMaxAge", 3600*12, `CSRF MaxAge`)
)

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
				<script src="static/debug.js?time=5"></script>
			</head>
			<body>
				Username: <input type="text" id="Username" value="jun" placeholder="Username" />
				<br>Password: <input type="text" id="Password" value="junpass" placeholder="Password" />
				<br>
				<br>User: <select id="User">
					<option value="user1">User1</option>
					<option value="user2">User2</option>
					<option value="user3">User3</option>
				</select>
				<br>CSRFToken: <input type="text" id="CSRFToken" style="width: 800px;" />
				<br>
				<br><button id="SalOrderBtn">SalOrder</button>
				<br>
				<br><button id="NewsBtn">News 1 (no CSRF)</button>
				<br><button id="News3Btn">News 3 (with CSRF)</button>
				<br><button id="CSRFBtn">Get CSRF Token</button>
				<br>
				<br><button id="LoginBtn">Login</button>
				<br><button id="LogoutBtn">Logout</button>
			</body>
		</html>
`
}

func srvHome(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	// Reference:
	// https://golang.org/pkg/net/http/#example_ServeMux_Handle
	// https://golang.org/pkg/net/http/#ServeMux
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if _, err := w.Write([]byte(homeHTML())); err != nil {
		log.Printf("w.Write: %v", err.Error())
	}
}

func srvLogin(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG_LOGIN: Inside")

	o := r.Context().Value("iojson").(*iojson.IOJSON)

	// User input
	i := iojson.NewIOJSON()

	if err := i.Decode(r.Body); err != nil {
		o.AddError(err.Error())
		return
	}

	log.Printf("DEBUG_LOGIN: Username: %s", i.GetData("SQLUsername"))
	log.Printf("DEBUG_LOGIN: Password: %s", i.GetData("SQLPassword"))
	log.Printf("DEBUG_LOGIN: JSON: %s", i.EncodePretty())

	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, *sessionName)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		o.AddError(err.Error())
		return
	}

	// TODO: add a logic to continue only when session.IsNew is true
	log.Printf("DEBUG_LOGIN: IsNew Session: %v", session.IsNew)

	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   *srvDomain,
		MaxAge:   *sessionMaxAge,
		Secure:   *sessionSecure, // TODO: set to true once applied the SSL certificate.
		HttpOnly: *sessionHttpOnly,
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

	o.AddData("good", "very good")
	o.AddData("username", username)
}

func srvLogout(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG_LOOUT: Inside")

	o := r.Context().Value("iojson").(*iojson.IOJSON)

	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, *sessionName)

	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//w.WriteHeader(http.StatusInternalServerError)
		o.AddError(err.Error())
		return
	}

	// TODO: add a logic to continue only when session.IsNew is false
	log.Printf("DEBUG_LOOUT: IsNew Session: %v", session.IsNew)

	log.Printf("DEBUG_LOOUT: Domain: %v", session.Options.Domain)
	log.Printf("DEBUG_LOOUT: MaxAge Before: %v", session.Options.MaxAge)

	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   *srvDomain,
		MaxAge:   -1,             // means delete cookie now.
		Secure:   *sessionSecure, // TODO: set to true once applied the SSL certificate.
		HttpOnly: *sessionHttpOnly,
	}

	log.Printf("DEBUG_LOOUT: Domain: %v", session.Options.Domain)
	log.Printf("DEBUG_LOOUT: MaxAge After: %v", session.Options.MaxAge)

	// Save it before we write to the response/return from the handler.
	if err := session.Save(r, w); err != nil {
		o.AddError(err.Error())
		return
	}

	log.Printf("DEBUG_LOOUT: ID: %v", session.ID)

	if err := os.Remove("./session/session_" + session.ID); err != nil {
		// TODO: do something
	}

	o.AddData("msg", "cookie has been deleted from server")
}

func srvSalOrder(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG_SalOrder: Inside")

	o := r.Context().Value("iojson").(*iojson.IOJSON)

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

	log.Printf("DEBUG_SALORDER: SQLDealerName: %s", i.GetData("SQLDealerName"))
	log.Printf("DEBUG_SALORDER: SQLIDShipAddr: %d", int(i.GetData("SQLIDShipAddr").(float64)))
	log.Printf("DEBUG_SALORDER: SQLPrice: %f", i.GetData("SQLPrice"))
	log.Printf("DEBUG_SALORDER: JSON: %s", i.EncodePretty())

	IDSalOrder := r.FormValue("IDSalOrder")

	o.AddData("IDSalOrder", IDSalOrder)
}

func srvNews1(w http.ResponseWriter, r *http.Request) {
	o := r.Context().Value("iojson").(*iojson.IOJSON)

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

func srvCSRFToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("DEBUG_CSRFToken: Inside")

	// Get the token and pass it in the CSRF header. Our JSON-speaking client
	// or JavaScript framework can now read the header and return the token in
	// in its own "X-CSRF-Token" request header on the subsequent POST.
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
}

//
func main() {
	//
	flag.Parse()

	//
	util.HelpGenTLSKeys()

	//
	homeChain := alice.New(
		gziphandler.GzipHandler,
		middleware.LoggerHandler,
		middleware.DomainHandler(*srvDomain+*srvPort),
	)

	minChain := alice.New(
		gziphandler.GzipHandler,
		iojson.EchoHandler,
		middleware.LoggerHandler,
		middleware.DomainHandler(*srvDomain+*srvPort),
	)

	stdChain := minChain.Append(
		middleware.AuthUserHandler(store, *sessionName),
		csrf.Protect([]byte(*csrfAuthKey),
			csrf.Domain(*srvDomain),
			csrf.Secure(*csrfSecure),
			csrf.MaxAge(*csrfMaxAge),
			csrf.ErrorHandler(iojson.ErrorHandler("Forbidden - CSRF token invalid")),
		),
	)

	//
	mux := http.NewServeMux()

	// NOTE: FileServer calls path.Clean() to clean up path.
	mux.Handle(*staticDir, gziphandler.GzipHandler(http.StripPrefix(*staticDir, http.FileServer(http.Dir("."+*staticDir)))))

	mux.Handle("/", homeChain.ThenFunc((srvHome)))

	mux.Handle("/Login", minChain.ThenFunc(srvLogin))
	mux.Handle("/Logout", minChain.ThenFunc(srvLogout))

	mux.Handle("/SalOrder", stdChain.ThenFunc(srvSalOrder))

	mux.Handle("/News1", minChain.ThenFunc(srvNews1))

	mux.Handle("/News3", stdChain.ThenFunc(srvNews1)) // with CSRF
	mux.Handle("/CSRFToken", stdChain.ThenFunc(srvCSRFToken))

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
