package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/voidt01/go-url-shortener/internal/models"
)

type Config struct {
	port     string
	baseURL  string
	dsn      string
	dbDriver string
}

type App struct {
	configApp *Config
	urls      *models.UrlStore
	errorLog  *log.Logger
	infoLog   *log.Logger
}

func main() {
	// app's config
	var config Config
	flag.StringVar(&config.port, "addr", ":4000", "HTTP Network Address Port")
	flag.StringVar(&config.baseURL, "base-url", "http://localhost", "Base URL for URL Shortener Service")
	flag.StringVar(&config.dsn, "dsn", "api:firstproject@/urlShortener", "Data Source Name for Database")
	flag.StringVar(&config.dbDriver, "driver", "mysql", "Database Driver name")
	flag.Parse()

	// app's logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := OpenDB(config.dbDriver, config.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Connected to database successfully")
	defer db.Close()

	app := &App{
		configApp: &config,
		urls:      models.NewStoreURL(),
		errorLog:  errorLog,
		infoLog:   infoLog,
	}

	srv := &http.Server{
		Addr:     config.port,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("starting a server on: %s", config.port)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func OpenDB(driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
