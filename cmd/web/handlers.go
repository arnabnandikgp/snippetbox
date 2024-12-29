package main

import (
	"fmt"
	"net/http"
	"strconv"
	"html/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	files := []string{ // directly fetches from the filesystem and not the file server(not to be confused)
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // will be logged by the custom logger
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil) // make the base html template visible on the route
	if err != nil {

		app.serverError(w, err) // will be logged by the custom logger
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	// w.Write([]byte("Hello from Snippetbox"))
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
