package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"chrisshyi.net/mini_url/pkg/models/postgres"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	miniURLModel *postgres.MiniURLModel
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	defaultDSN := "postgres://mini_url:mini_pass@localhost:5432/mini_url"
	dsn := flag.String("dsn", defaultDSN, "The database connection string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		miniURLModel: &postgres.MiniURLModel{},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
