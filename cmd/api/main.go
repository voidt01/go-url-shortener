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
	Server struct {
		port    string
		baseURL string
	}
	Database struct {
		dsn      string
		dbDriver string
	}
}

type App struct {
	configApp *Config
	urls      *models.UrlModel
	errorLog  *log.Logger
	infoLog   *log.Logger
}

func main() {
	// app's config
	var cfg Config
	flag.StringVar(&cfg.Server.port, "addr", ":4000", "HTTP Network Address Port")
	flag.StringVar(&cfg.Server.baseURL, "base-url", "http://localhost", "Base URL for URL Shortener Service")
	flag.StringVar(&cfg.Database.dsn, "dsn", "api:firstproject@/urlShortener?parseTime=true", "Data Source Name for Database")
	flag.StringVar(&cfg.Database.dbDriver, "driver", "mysql", "Database Driver name")
	flag.Parse()

	// app's logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Pool of DB Conn
	db, err := OpenDB(cfg.Database.dbDriver, cfg.Database.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Connected to database successfully")
	defer db.Close()

	app := &App{
		configApp: &cfg,
		urls:      &models.UrlModel{DB: db},
		errorLog:  errorLog,
		infoLog:   infoLog,
	}

	srv := &http.Server{
		Addr:     cfg.Server.port,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("starting a server on %s", cfg.Server.port)
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
