package main

import (
	"database/sql"
	"flag"
	"github.com/arnabnandikgp/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
)

// application struct
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {

	//  making custom loggers in order to make the logs more readable and constant in whole project.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	//  flags return pointers to the flag values
	//  addr flag is passed as an argument when running the server
	addr := flag.String("addr", ":4000", "HTTP network address")
	//  command line flag for the database DSN(mysql DSN(data source name))
	dsn := flag.String("dsn", "web:20016@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	// the *sql.DB object is stored in db
	db, err := openDB(*dsn)

	// creating a new application struct to make the custom loggers available to the handlers
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	//  info logger
	infoLog.Printf("Starting server on %s", *addr)

	//  creating a new http.ServeMux defined in the routes.go
	mux := app.routes()

	if err != nil {
		errorLog.Fatal(err) //fatal error if database can't be opened
	}

	defer db.Close()

	//  custom http.Server for making the custom loggers available
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Call the ListenAndServe() method on our new http.Server struct instead of  err := http.ListenAndServe(*addr, mux)
	err = srv.ListenAndServe()

	//  error logger
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
