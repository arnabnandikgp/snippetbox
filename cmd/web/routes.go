package main

import (
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
    fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
    f, err := nfs.fs.Open(path)
    if err != nil {
        return nil, err
    }

    s, err := f.Stat()
    if err != nil {
        return nil, err
    }
    
    if s.IsDir() {
        index := filepath.Join(path, "index.html")
        if _, err := nfs.fs.Open(index); err != nil {
            closeErr := f.Close()
            if closeErr != nil {
                return nil, closeErr
            }
            return nil, err
        }
    }

    return f, nil
}    
// The routes() method returns a servemux containing our application routes.
func (app *application) routes() *http.ServeMux {
	
mux := http.NewServeMux()

//  passing the handler functions as methods of the application struct
mux.HandleFunc("/", app.home)
mux.HandleFunc("/snippet/view", app.snippetView)
mux.HandleFunc("/snippet/create", app.snippetCreate)

//  file server added and allows access for directories with html files.
fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
mux.Handle("/static", http.NotFoundHandler())
mux.Handle("/static/", http.StripPrefix("/static", fileServer))

return mux
}