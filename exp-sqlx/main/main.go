package main

import (
	"fmt"
	"net/http"
)

import (
	"github.com/junhsieh/goexamples/exp-sqlx/model"
	"github.com/junhsieh/goexamples/exp-sqlx/router"
)

func main() {
	var err error

	if err = model.Open(); err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}

	defer model.Close()

	//
	router.RegisterRoute()

	//
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
		return
	}
}
