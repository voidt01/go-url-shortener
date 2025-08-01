package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/voidt01/go-url-shortener/internal/url"
	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	config   *Config
	errorLog *log.Logger
	infoLog  *log.Logger
	URLHandler *url.UrlHandler
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	cfg, err := LoadConfig()
	if err != nil {
		errorLog.Fatal(err)
	}

	db, err := OpenDB(cfg.database.dbDriver, cfg.database.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Connected to database successfully")
	defer db.Close()

	URLStore := url.NewUrlStore(db)
	URLService := url.NewUrlService(URLStore)

	app := &Application{
		config:   cfg,
		errorLog: errorLog,
		infoLog:  infoLog,
		URLHandler: url.NewUrlHandler(URLService),
	}

	err = app.serve()
	if err != nil {
		errorLog.Fatal(err)
	}
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

