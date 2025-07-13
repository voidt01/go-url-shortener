package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	config   *Config
	errorLog *log.Logger
	infoLog  *log.Logger
	URLHandler *urlHandler
}

func main() {
	// app's logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// app's config
	cfg, err := LoadConfig()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Pool of DB Conn
	db, err := OpenDB(cfg.database.dbDriver, cfg.database.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Connected to database successfully")
	defer db.Close()

	URLStore := NewUrlStore(db)
	URLService := NewUrlService(URLStore)

	app := &Application{
		config:   cfg,
		errorLog: errorLog,
		infoLog:  infoLog,
		URLHandler: NewUrlHandler(URLService),
	}

	srv := &http.Server{
		Addr:     cfg.port,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("starting a server on %s", cfg.port)
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
