package main

//go:generate go run gen.go
//go:generate go build -o main vfsgen_static_assets.go assets_vfsdata.go

import (
	"github.com/gin-gonic/gin"
)

func main() {
	//
	r := gin.New()

	r1 := r.Group("/")
	r1.Use()
	{
		// NOTE: curl localhost:8080/dist/dist/app.js
		r1.StaticFS("/dist", assets)
	}

	r.Run()
}
