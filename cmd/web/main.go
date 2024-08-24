package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"golangify.com/snippetbox/pkg/models/pgsql"
)




type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *pgsql.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn","postgres://postgres:BUGLb048@localhost:5432/test", "Название PgSQL источника данных")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()



	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &pgsql.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	//подключение к бд
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	// Проверка соединения с базой данных
    err = db.Ping(context.Background())
    if err != nil {
        db.Close()
        return nil, err
    }
    
	return db, nil
}