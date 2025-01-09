package main

import (
	"fmt"
	"net/http"
	"strconv"
	"html/template"
	"github.com/arnabnandikgp/snippetbox/internal/models"
	"errors"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%+v\n", snippet)
	// }
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

	data := &templateData{
		Snippets: snippets,
	}

	err = ts.ExecuteTemplate(w, "base",data) // make the base html template visible on the route
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

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return 
	}

	data := &templateData{
		Snippet : snippet,
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
	}
	// Parse the template files...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
	// fmt.Fprintf(w, "%+v", snippet)
	// fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content :="O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– KobayashiIssa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)

	// w.Write([]byte("Create a new snippet..."))
}
// new sys