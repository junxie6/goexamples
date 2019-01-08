package main

import (
	"log"
	"net/http"
)
import (
	"github.com/shurcooL/httpfs/union"
	"github.com/shurcooL/vfsgen"
)

func main() {
	var err error

	var Assets = union.New(map[string]http.FileSystem{
		"/dist": http.Dir("./dist"),
	})

	if err = vfsgen.Generate(Assets, vfsgen.Options{}); err != nil {
		log.Fatalln(err)
	}
}
