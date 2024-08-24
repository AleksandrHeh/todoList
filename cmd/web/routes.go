package main

import "net/http"

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func (app *application) routes() *http.ServeMux{
	// Используем методы из структуры в качестве обработчиков маршрутов.
	mux := http.NewServeMux()

	
	mux.Handle("/api/register", enableCORS(http.HandlerFunc(app.register)))
	mux.HandleFunc("/api/login", app.login)

	return mux
}	