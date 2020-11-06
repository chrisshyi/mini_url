package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"chrisshyi.net/mini_url/pkg/models"
	"chrisshyi.net/mini_url/pkg/models/postgres"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	miniURLModel interface {
		GetByID(int) (*models.MiniURL, error)
		GetByURL(string) (*models.MiniURL, error)
		Insert(string) (int, error)
	}
}

type config struct {
	addr   *string
	dbHost *string
	dbPort *string
	dbUser *string
	dbPass *string
	dbName *string
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

func readEnvOrDefault(key, defaultVal string) *string {
	envVal, exists := os.LookupEnv(key)
	if !exists {
		return &defaultVal
	}
	return &envVal
}

func appConfigSetup() config {
	appConfig := config{
		addr:   readEnvOrDefault("APP_ADDR", ":4000"),
		dbHost: readEnvOrDefault("DB_HOST", "localhost"),
		dbPort: readEnvOrDefault("DB_PORT", "5432"),
		dbUser: readEnvOrDefault("DB_USER", "mini_url"),
		dbPass: readEnvOrDefault("DB_PASS", "mini_pass"),
		dbName: readEnvOrDefault("DB_NAME", "mini_url"),
	}
	return appConfig
}

// @title URL shortening API
// @version 1.0
// @description This is a URL shortening service
// @host localhost:4000
// @BasePath /
func main() {
	appConfig := appConfigSetup()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		*appConfig.dbUser,
		*appConfig.dbPass,
		*appConfig.dbHost,
		*appConfig.dbPort,
		*appConfig.dbName,
	)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		miniURLModel: &postgres.MiniURLModel{DB: db},
	}

	srv := &http.Server{
		Addr:         *appConfig.addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *appConfig.addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
