package main

import (

	"errors"
	"fmt"
	"html/template"
	_ "html/template"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"golangify.com/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
        app.notFound(w)
        return
    }
 
    s, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, err)
        return
    }
 
    for _, snippet := range s {
        fmt.Fprintf(w, "%v\n", snippet)
    }

}

func (app *application) showSnippetTest(w http.ResponseWriter, r *http.Request){
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Страница не найдена.
		return
	}
 
	// Вызываем метода Get из модели Snipping для извлечения данных для
	// конкретной записи на основе её ID. Если подходящей записи не найдено,
	// то возвращается ответ 404 Not Found (Страница не найдена).
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
 
	// Отображаем весь вывод на странице.
	fmt.Fprintf(w, "%v", s)
}
 
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    key := fmt.Sprintf("snippet:%d", id)
    var s *models.Snippet

    // Попробуйте получить заметку из кэша Redis
    cachedNote, err := app.redis.Get(key).Result()
    if err == redis.Nil {
        // Если заметка не найдена в кэше, получите её из базы данных
        s, err = app.snippets.Get(id)
        if err != nil {
            if errors.Is(err, models.ErrNoRecord) {
                app.notFound(w)
            } else {
                app.serverError(w, err)
            }
            return
        }

        // Сохраните заметку в кэш Redis
        err = app.redis.Set( key, s.Content, 0).Err()
        if err != nil {
            app.errorLog.Printf("Ошибка кэширования заметки: %v", err)
        }
    } else if err != nil {
        // Ошибка при получении заметки из Redis
        app.serverError(w, err)
        return
    } else {
        // Если данные найдены в кэше, преобразуйте их обратно в структуру Snippet
        s = &models.Snippet{Content: cachedNote}
    }

    // Создаем экземпляр структуры templateData, содержащей данные заметки.
    data := &templateData{Snippet: s}

    files := []string{
        "./ui/html/show.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
        app.serverError(w, err)
        return
    }

    // Передаем структуру templateData в качестве данных для шаблона.
    err = ts.Execute(w, data)
    if err != nil {
        app.serverError(w, err)
    }
}
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed) // Используем помощник clientError()
		return
	}

	// Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.
	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"

		// Передаем данные в метод SnippetModel.Insert(), получая обратно
	// ID только что созданной записи в базу данных.

	id, err := app.snippets.Insert(title,content,expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Перенаправляем пользователя на соответствующую страницу заметки.
	http.Redirect(w,r,fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func dowlandHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w,r, "./ui/static/file.zip")
}