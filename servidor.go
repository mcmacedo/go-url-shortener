package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	porta   int
	urlBase string
)

func init() {
	porta = 8888
	urlBase = fmt.Sprintf("http://locahost:%d", porta)
}

func main() {
	http.HandleFunc("/api/encurtar", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "%s", "encurtar")
	})
	http.HandleFunc("/r/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "%s", "redirecionar")
	})

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", porta), nil))
}
