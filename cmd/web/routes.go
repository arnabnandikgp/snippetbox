package main

import (
	"net/http"
	"path/filepath"
    "github.com/justinas/alice"
    "github.com/julienschmidt/httprouter"
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
func (app *application) routes() http.Handler {
	
    // mux := http.NewServeMux()
    router := httprouter.New()

    router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.notFound(w)
    })

    //  file server added and allows access for directories with html files.
    fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
    // mux.Handle("/static", http.NotFoundHandler())
    // mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    //  passing the handler functions as methods of the application struct
    router.HandlerFunc(http.MethodGet, "/", app.home)
    router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
    router.HandlerFunc(http.MethodGet, "/snippet/createForm", app.snippetCreateForm)
    router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreate)


    standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

    return standard.Then(router)
    // return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}