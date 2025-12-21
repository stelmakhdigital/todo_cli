package render

import (
	"fmt"
	"strings"
	"todo_cli/internal/task"
	"unicode/utf8"
)

type Render interface {
	RenderList(tasks []*task.Task)
	RenderMap(data map[string]interface{})
	RenderDetailed(tasks *task.Task)
}

type TerminalRender struct{}

func (r *TerminalRender) RenderList(tasks []*task.Task) {
	var columnMax int = 0
	for _, value := range tasks {
		if columnMax < utf8.RuneCountInString(value.Title) {
			columnMax = utf8.RuneCountInString(value.Title)
		}
	}
	columnMax += 5
	fmt.Print("\n")
	fmt.Printf("%-4s | %-*s | %-12s | %-15s\n", "ID", columnMax, "Название", "Статус", "Создана")
	fmt.Println(strings.Repeat("-", columnMax+40))

	for _, task := range tasks {
		fmt.Printf("%-4d | %-*s | %-12s | %-15s\n",
			task.ID, columnMax, task.Title, task.Status,
			task.CreatedAt.Format("02.01.2006"))
	}
	fmt.Print("\n")
}

func (r *TerminalRender) RenderMap(data map[string]interface{}) {
	fmt.Print("\n")
	for name, value := range data {
		if stringValue, ok := value.(string); ok {
			fmt.Printf("%v: %v\n", name, stringValue)
		}
		if intValue, ok := value.(int); ok {
			fmt.Printf("%v: %v\n", name, intValue)
		}
	}
	fmt.Print("\n")
}

func (r *TerminalRender) RenderDetailed(tasks *task.Task) {
	fmt.Print("\n")
	fmt.Printf("ID: %d\n", tasks.ID)
	fmt.Printf("Название: %s\n", tasks.Title)
	fmt.Printf("Описание: %s\n", tasks.Description)
	fmt.Printf("Статус: %s\n", tasks.Status.String())
	fmt.Printf("Создана: %s\n", tasks.CreatedAt.Format("02.01.2006 15:04"))
	fmt.Print("\n")
}
