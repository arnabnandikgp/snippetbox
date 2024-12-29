package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

//  application struct
type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	}

func main() {

	//  making custom loggers in order to make the logs more readable and constant in whole project.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	// creating a new application struct to make the custom loggers available to the handlers
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		}
	
	//  addr flag is passed as an argument when running the server	
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	//  info logger
    infoLog.Printf("Starting server on %s", *addr)

	//  creating a new http.ServeMux defined in the routes.go
	mux := app.routes()
	
	//  custom http.Server for making the custom loggers available
 	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
		}

	// Call the ListenAndServe() method on our new http.Server struct instead of  err := http.ListenAndServe(*addr, mux)
	err := srv.ListenAndServe()

	//  error logger
	errorLog.Fatal(err)
}
