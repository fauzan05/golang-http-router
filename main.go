package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var Localhost string = "localhost:8000"
var FullLocalhost string = "http://localhost:8000/"
func main(){
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Fprint(w, "testing http router")
	})

	server := http.Server{
		Handler: router,
		Addr: Localhost,
	}

	server.ListenAndServe()
}