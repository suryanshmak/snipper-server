package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"suryanshmak.net/snippetBox/pkg/models/postgres"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *postgres.SnippetModel
	users    *postgres.UserModel
	// templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "postgres://web:pass1234@localhost:5432/snippetbox", "PosgreSQL data source name")
	flag.Parse()
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close(context.Background())

	// templateCache, err := newTemplateCache("./ui/html")
	// if err != nil {
	// 	errorLog.Fatal(err)
	// }

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &postgres.SnippetModel{DB: db},
		users:    &postgres.UserModel{DB: db},
		// templateCache: templateCache,
	}

	if err != nil {
		errorLog.Fatal(err)
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}
	return db, nil
}
