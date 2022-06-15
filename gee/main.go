package main

import (
	"fmt"
	"gee2"
	"net/http"
)

func main() {
	r := gee2.New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL)
	})
	r.Run(":9999")
}
