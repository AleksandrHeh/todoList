package main

import (
	"encoding/json"
	_ "errors"
	"fmt"
	_ "fmt"
	"net/http"
	_ "golangify.com/snippetbox/pkg/models"
)

type User struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
	MiddleName string `json:"middlename"`
    Email    string `json:"email"`
    Password string `json:"password"`
}


func (app *application) register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        fmt.Println("Error decoding request body:", err)
        return
    }

	exists, err := app.snippets.EmailExists(user.Email)
	if err != nil {
		http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
        return
	}

	if exists{
		http.Error(w,"Email already in use", http.StatusConflict)
		return
	}

    id, err := app.snippets.InsertUser(user.FirstName, user.LastName, user.MiddleName, user.Email, user.Password)
    if err != nil {
        http.Error(w, "Unable to create user", http.StatusInternalServerError)
        fmt.Println("Error inserting user:", err)
        return
    }

	response := map[string]interface{}{
        "id": id,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}


func (app *application) login(w http.ResponseWriter, r *http.Request) {
    
}


/*func (app *application) showSnippetTest(w http.ResponseWriter, r *http.Request) {
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

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		
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
		err = app.redis.Set(key, s.Content, 0).Err()
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

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Перенаправляем пользователя на соответствующую страницу заметки.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func dowlandHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/static/file.zip")
}*/