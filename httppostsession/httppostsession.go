package main

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

//var store = sessions.NewCookieStore([]byte("something-very-secret"))
var store = sessions.NewFilesystemStore("./session", []byte("something-very-secret"))

func handler(w http.ResponseWriter, r *http.Request) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "MySession")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   "erp.local",
		MaxAge:   86400 * 7, // 7 Days
		Secure:   false,
		HttpOnly: true,
	}

	if r.FormValue("act") == "set" {
		// Set some session values.
		session.Values["name"] = r.FormValue("name")

		// Save it before we write to the response/return from the handler.
		session.Save(r, w)
	}

	log.Printf("Host: %v", r.Host)
	log.Printf("act: %v; session: %v", r.FormValue("act"), session)

	name := "init name"

	if s, ok := session.Values["name"]; ok {
		name = s.(string)
	}

	for k, v := range session.Values {
		log.Printf("%v: %v", k, v)
	}
	log.Printf("====================")

	w.Write([]byte(getHTML("Hi %s", name)))
}

func getHTML(format string, v ...interface{}) string {
	return fmt.Sprintf(`
		<html>
			<head>
				<link rel="shortcut icon" href="data:image/x-icon;," type="image/x-icon">
			</head>
			<body>`+format+
		`	</body>
		</html>`, v...)
}

func main() {
	http.HandleFunc("/", handler)

	// Important Note: If you aren't using gorilla/mux, you need to wrap your handlers with
	// context.ClearHandler as or else you will leak memory!
	// An easy way to do this is to wrap the top-level mux when calling http.ListenAndServe:
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
