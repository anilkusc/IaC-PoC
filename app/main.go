package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	message := "Hello :" + string(r.URL.Query().Get("Name"))
	io.WriteString(w, message)
	return

}

func main() {
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/", Hello)
	fmt.Println("Serving on:8080")
	http.ListenAndServe(":8080", r)

}
