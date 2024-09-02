package main

import (
	"log"
	"net/http"
)

// enableCORS добавляет заголовки CORS ко всем запросам
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Handling request: %s %s", r.Method, r.URL.Path)

        w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешаем запросы с любого домена
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Если это OPTIONS-запрос, просто возвращаем успешный ответ
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Продолжаем обработку запроса
        next.ServeHTTP(w, r)
    })
}

func (app *application) routes() *http.ServeMux {
    mux := http.NewServeMux()

    // Регистрация обработчиков с использованием http.Handler
    mux.Handle("/api/register", enableCORS(http.HandlerFunc(app.register)))
    mux.Handle("/api/login", enableCORS(http.HandlerFunc(app.login)))
    mux.Handle("/api/board/createProject", enableCORS(http.HandlerFunc(app.createProject)))
    mux.Handle("/api/board/userCreatedProjects", enableCORS(http.HandlerFunc(app.userCreatedProjects)))

    return mux
}