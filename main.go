package main

import (
	"log"
	"net/http"
	)
	// Define a home handler function which writes a byte slice containing
	// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
		}
	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snippetview handler"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from snippetCreate handler"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home) 
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux) 
	log.Fatal(err)
}