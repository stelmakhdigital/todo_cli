package main

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidStatus = errors.New("некорректный статус")
	ErrInvalidID     = errors.New("некорректный ID")
	ErrTaskNotFound  = errors.New("задача не найдена")
)

// type-alias
type Status string

const (
	StatusPending   Status = "pending"
	StatusProgress  Status = "in_progress"
	StatusCompleted Status = "completed"
)

// теги структур json или yaml задаются через тильда кавычки
// omitempty - пропустить если нету значения
// "-" - исключить вообще
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created,omitempty"`
	CompletedAt *time.Time `json:"completed,omitempty"`
}

// это метод - функция с получателем (receiver)
// реализует интерфейс для вывода значения с преобразованием типов
// fmt.PrintLn(StatusPending) вызовет StatusPending.String() и выведет "pending"
func (status Status) String() string {
	return string(status)
}

func (status Status) Valid() bool {
	switch status {
	case StatusPending, StatusProgress, StatusCompleted:
		return true
	}
	return false
}

// конструктор указывается через New префикс. У конструктора нет именованных аргументов
// %w позволяет обернуть ошибку для error.Is() проверки
func NewTask(id int, title, description, status string) (Task, error) {
	if !Status(status).Valid() {
		return Task{}, fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}
	return Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      Status(status),
		CreatedAt:   time.Now(),
	}, nil
}

func main() {
	fmt.Println("Hello")
	fmt.Println("Time: ", time.Now().Format("02.01.2006 15:05"))
	tsk, err := NewTask(1, "test", "dest", "pending")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tsk)
}
