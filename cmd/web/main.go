package main

import (
	"flag"
	"log"
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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	// fileServer := http.FileServer(http.Dir("./ui/static/")) //handle static files
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
    mux.Handle("/static", http.NotFoundHandler())
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
