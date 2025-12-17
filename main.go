package main

import (
	"todo_cli/internal/manager"
	"todo_cli/internal/storage"
)

func main() {
	store := &storage.FileStorage{}
	filter := &manager.FilterTasks{}
	mgr := manager.NewManager(store, filter)

	// allTasks, err := mgr.GetAll()
	// if err != nil {
	// 	log.Error(err)
	// }
	// mgr.Render(allTasks)

	// testTask := map[string]string{
	// 	"Title":       "Test task",
	// 	"Description": "bla vla bla",
	// }

	// mgr.Create(testTask)

	// mgr.Start(8)

	// allTasks2, err2 := mgr.GetAll()
	// if err2 != nil {
	// 	log.Error(err2)
	// }
	mgr.Show(10)

}
