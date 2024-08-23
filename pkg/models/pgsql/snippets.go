package pgsql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golangify.com/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashPassword,err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPassword),err
}

func checkHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password),[]byte(hash))
	return err == nil
}


// SnippetModel - Определяем тип который обертывает пул подключения sql.DB
type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) insertUser(email, password string)(int,error) {
	hashPassword, err := hashPassword(password)
	stmt := "INSERT INTO snippet (email, password) VALUES ($1,$2)"
}

// Insert - Метод для создания новой заметки в базе дынных.
/*func (m *SnippetModel) Insert (title, content, expires string) (*models.User,error){
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES ($1, $2, NOW(), NOW() + $3::interval) RETURNING id`

	var id int
	err :=  m.DB.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)
	if err != nil{
		return 0, nil
	}

	return int(id), nil
}*/

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *SnippetModel) Get (id int) (*models.Snippet, error){
	
	stmt :=  "SELECT id, title, content, created, expires FROM snippets WHERE id = $1 "

	
	row := m.DB.QueryRow(context.Background(),stmt, id)

	// Инициализируем указатель на новую структуру Snippet.
	s := &models.Snippet{}

	// Используйте row.Scan(), чтобы скопировать значения из каждого поля от sql.Row в 
	// соответствующее поле в структуре Snippet. Обратите внимание, что аргументы 
	// для row.Scan - это указатели на место, куда требуется скопировать данные
	// и количество аргументов должно быть точно таким же, как количество 
	// столбцов в таблице базы данных.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// Специально для этого случая, мы проверим при помощи функции errors.Is()
		// если запрос был выполнен с ошибкой. Если ошибка обнаружена, то
		// возвращаем нашу ошибку из модели models.ErrNoRecord.
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else{
			return nil, err
		}
	}

	// Если все хорошо, возвращается объект Snippet
	return s,nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > CURRENT_TIMESTAMP ORDER BY created DESC LIMIT 10 "

	// Используем метод Query() для выполнения нашего SQL запроса.
	// В ответ мы получим sql.Rows, который содержит результат нашего запроса.
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}

	// правильно закроется перед вызовом метода Latest(). Этот оператор откладывания
	// должен выполнится *после* проверки на наличие ошибки в методе Query().
	// В противном случае, если Query() вернет ошибку, это приведет к панике
	// так как он попытается закрыть набор результатов у которого значение: nil.
	defer rows.Close()

	// Инициализируем пустой срез для хранения объектов models.Snippets.
	var snippets []*models.Snippet

	// Используем rows.Next() для перебора результата. Этот метод предоставляем
	// первый а затем каждую следующею запись из базы данных для обработки
	// методом rows.Scan().
	for rows.Next() {
		// Создаем указатель на новую структуру Snippet
		s := &models.Snippet{}
		// Используем rows.Scan(), чтобы скопировать значения полей в структуру.
		// Опять же, аргументы предоставленные в row.Scan()
		// должны быть указателями на место, куда требуется скопировать данные и
		// количество аргументов должно быть точно таким же, как количество
		// столбцов из таблицы базы данных, возвращаемых вашим SQL запросом.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
				// Добавляем структуру в срез.
				snippets = append(snippets, s)
	}
 

	// Когда цикл rows.Next() завершается, вызываем метод rows.Err(), чтобы узнать
	// если в ходе работы у нас не возникла какая либо ошибка.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}