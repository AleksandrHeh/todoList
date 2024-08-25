package main

import (
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Handling request: %s %s", r.Method, r.URL.Path)

        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}


func (app *application) routes() *http.ServeMux {
    mux := http.NewServeMux()

    // Регистрация обработчиков с использованием http.Handler
    mux.Handle("/api/register", enableCORS(http.HandlerFunc(app.register)))
    mux.Handle("/api/login", enableCORS(http.HandlerFunc(app.login)))

    return mux
}
