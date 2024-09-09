package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено!")

type User struct{
	UserID int
	FirstName string
	LastName string
	MiddleName string
	Email string
	Password string
}

type Project struct{
	ProjectID int
	ProjectName string
	Description string
	Password string
	CreatedBy int //id users кто создал проект
}

