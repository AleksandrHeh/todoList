package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golangify.com/snippetbox/pkg/models/pgsql"
)

type user struct {
	FirstName 	string `json:"firstname"`
	LastName 	string `json:"lastname"`
	MiddleName 	string `json:"middlename"`
    Email    	string `json:"email"`
    Password 	string `json:"password"`
}

type project struct{
    ProjectName   string `json:"projectname"`
    Description   string `json:"description"`
    Password      string `json:"password"`
    Email string `json:"email"` // Должно совпадать с ключом в JSON
}


// Структура для хранения заявлений JWT
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (app *application) createProject(w http.ResponseWriter, r *http.Request){
    if r.Method != http.MethodPost{
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var p project

    err := json.NewDecoder(r.Body).Decode(&p)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    log.Printf("Данные проекта: %+v", p)
    
    //Получаем id пользователя по email
    userID,err := app.snippets.GetUserByEmail(p.Email)
    if err != nil {
        log.Printf("Error retrieving user ID: %v", err)
    http.Error(w, "Error retrieving user ID", http.StatusInternalServerError)
    return
    }

    projectID, err := app.snippets.InserProject(p.ProjectName, p.Description, p.Password, userID) 
    if err != nil {
        if err != nil {
            log.Printf("Error creating project: %v", err)
            http.Error(w, "Error creating project", http.StatusInternalServerError)
            return
        }
    }

    response := map[string]interface{}{
        "success": true,
        "projectID": projectID,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

var jwtKey = []byte("your_super_secret_key_which_is_long_and_random_enough")

// Функция для создания JWT токена
func createToken (email string) (string, error) {
	expiratonTime := time.Now().Add(24*time.Hour)

	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiratonTime),
		},

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSrting, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenSrting, nil
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var user user
	//Декодируем запрос
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
    var user user

    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    log.Printf("Received login request for email: %s", user.Email)

    exists, err := app.snippets.GetUserAuthorization(user.Email)
    if err != nil {
        log.Printf("Error retrieving user: %v", err)
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    hash, err := pgsql.HashPassword(user.Password)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    if hash != exists.Password {
        log.Printf("Password mismatch for email: %s", user.Email)
        http.Error(w, "Invalid email or password", http.StatusUnauthorized)
        return
    }

    // Создаем JWT Token для авторизованного пользователя
    token, err := createToken(user.Email)
    if err != nil {
        http.Error(w, "Could not create token", http.StatusInternalServerError)
		log.Printf("Error signing token: %v", err) // Логирование ошибки
        return
    }

    response := map[string]interface{}{
        "token": token,
        "email": user.Email,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
