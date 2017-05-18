package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: this test is not done yet. Leave it for example
func TestPublishWrongResponseStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))

	defer ts.Close()

	postURL := ts.URL

	err := example2(postURL)

	if err == nil {
		t.Errorf("Publish() didn’t return an error")
	}
}

// TODO: this test is not done yet. Leave it for example
func TestPublishWrongResponseStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})

	srv := &http.Server{
		Handler: mux,
		//Addr:    srvAddr,
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout:   srvWriteTimeout,
		//ReadTimeout:    srvReadTimeout,
		//MaxHeaderBytes: srvMaxHeaderBytes,
	}

	ts := httptest.NewUnstartedServer(nil)
	ts.Config = srv
	ts.Start()

	defer ts.Close()

	postURL := ts.URL

	err := example2(postURL)

	if err == nil {
		t.Errorf("Publish() didn’t return an error")
	}
}
