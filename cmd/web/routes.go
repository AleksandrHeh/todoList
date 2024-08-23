package main

import "net/http"

func (app *application) routes() *http.ServeMux{
	// Используем методы из структуры в качестве обработчиков маршрутов.
	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", app.login)
    mux.HandleFunc("/api/register", app.register)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}	