package main

import (
	"todo_cli/internal/manager"

	"github.com/labstack/gommon/log"
)

func main() {
	// var slice_t []*task.Task
	// // Создаём задачу через конструктор
	// t1, err := task.NewTask(1, "Купить молоко", "Зайти в магазин", string(task.StatusPending))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// slice_t = append(slice_t, t1)

	// t2, err := task.NewTask(2, "Купить молоко 2 - тестироваине строки на максимальное кол-во символов", "Зайти в магазин", string(task.StatusCompleted))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// slice_t = append(slice_t, t2)

	// t3, _ := task.NewTask(3, "Купить молоко 3", "Зайти в магазин 3", string(task.StatusPending))

	// slice_t = append(slice_t, t3)

	// taskToString, _ := storage.DataToJson(&slice_t)
	// storage.Save(taskToString)

	// stringToTasks, _ := storage.Load()

	// var task []*task.Task
	// storage.JsonToData(stringToTasks, &task)
	// fmt.Printf("String %v", task)

	// fmt.Printf("%-4s | %-70s | %-12s | %-15s\n", "ID", "Название", "Статус", "Создана")
	// fmt.Println(strings.Repeat("-", 111))
	// for _, value := range task {
	// 	// fmt.Printf("Задача #%d: %s (статус: %s)\n", value.ID, value.Title, value.Status)

	// 	// %-4s - строка, выравнивание влево, ширина 4 символа
	// 	// %-20s - строка, ширина 20 символов
	// 	// %d - число, %-12s - строка

	// 	fmt.Printf("%-4d | %-70s | %-12s | %-15s\n", value.ID, value.Title, value.Status, value.CreatedAt.Format("02.01.2006"))
	// }

	allTasks, err := manager.GetAll()
	if err != nil {
		log.Error(err)
	}
	manager.Render(allTasks)
}
