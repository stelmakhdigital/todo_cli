package main

import (
	"testing"
	"todo_cli/internal/manager"
	"todo_cli/internal/task"
)

type MockStorage struct {
	tasks []*task.Task
}

func (m *MockStorage) Load(fileName *string) ([]*task.Task, error) {
	return m.tasks, nil
}

func (m *MockStorage) Save(tasks []*task.Task, fileName *string) error {
	m.tasks = tasks
	return nil
}

func (m *MockStorage) Clear(fileName *string) error {
	return nil
}

type MockRender struct{}

func (r *MockRender) RenderList(tasks []*task.Task)         {}
func (r *MockRender) RenderMap(data map[string]interface{}) {}
func (r *MockRender) RenderDetailed(tasks *task.Task)       {}

func BenchmarkCreateTasks(b *testing.B) {
	store := &MockStorage{tasks: []*task.Task{}}
	filter := &manager.FilterTasks{}
	render := &MockRender{}
	mgr := manager.NewManager(store, filter, render)
	testData := map[string]string{
		"title":       "Тест",
		"description": "Тест",
	}
	b.ResetTimer()
	for b.Loop() {
		mgr.Create(testData)
	}
}

func BenchmarkEditTasks(b *testing.B) {
	store := &MockStorage{tasks: []*task.Task{}}
	filter := &manager.FilterTasks{}
	render := &MockRender{}
	mgr := manager.NewManager(store, filter, render)
	testData := map[string]string{
		"title":       "Тест",
		"description": "Тест",
	}
	newData := map[string]string{
		"title":       "Тест2",
		"description": "Тест2",
	}
	id, _ := mgr.Create(testData)
	b.ResetTimer()
	for b.Loop() {
		mgr.Edit(*id, newData)
	}
}

func BenchmarkStartTasks(b *testing.B) {
	store := &MockStorage{tasks: []*task.Task{}}
	filter := &manager.FilterTasks{}
	render := &MockRender{}
	mgr := manager.NewManager(store, filter, render)
	testData := map[string]string{
		"title":       "Тест",
		"description": "Тест",
	}
	id, _ := mgr.Create(testData)
	b.ResetTimer()
	for b.Loop() {
		mgr.Start(*id)
	}
}

func BenchmarkCompletetTasks(b *testing.B) {
	store := &MockStorage{tasks: []*task.Task{}}
	filter := &manager.FilterTasks{}
	render := &MockRender{}
	mgr := manager.NewManager(store, filter, render)
	testData := map[string]string{
		"title":       "Тест",
		"description": "Тест",
	}
	id, _ := mgr.Create(testData)
	b.ResetTimer()
	for b.Loop() {
		mgr.Complete(*id)
	}
}

func BenchmarkDeletetTasks(b *testing.B) {
	store := &MockStorage{tasks: []*task.Task{}}
	filter := &manager.FilterTasks{}
	render := &MockRender{}
	mgr := manager.NewManager(store, filter, render)
	testData := map[string]string{
		"title":       "Тест",
		"description": "Тест",
	}
	id, _ := mgr.Create(testData)
	b.ResetTimer()
	for b.Loop() {
		mgr.Delete(*id)
	}
}

func BenchmarkShowTasks(b *testing.B) {
	store := &MockStorage{tasks: []*task.Task{}}
	filter := &manager.FilterTasks{}
	render := &MockRender{}
	mgr := manager.NewManager(store, filter, render)
	testData := map[string]string{
		"title":       "Тест",
		"description": "Тест",
	}
	id, _ := mgr.Create(testData)
	b.ResetTimer()
	for b.Loop() {
		mgr.Show(*id)
	}
}

func BenchmarkListTasks(b *testing.B) {
	store := &MockStorage{tasks: []*task.Task{}}
	filter := &manager.FilterTasks{}
	render := &MockRender{}
	mgr := manager.NewManager(store, filter, render)
	testData := map[string]string{
		"title":       "Тест",
		"description": "Тест",
	}
	mgr.Create(testData)
	b.ResetTimer()
	for b.Loop() {
		mgr.List("all")
	}
}

func BenchmarkSearchTasks(b *testing.B) {
	store := &MockStorage{tasks: []*task.Task{}}
	filter := &manager.FilterTasks{}
	render := &MockRender{}
	mgr := manager.NewManager(store, filter, render)
	testData := map[string]string{
		"title":       "Тест",
		"description": "Тест",
	}
	mgr.Create(testData)
	b.ResetTimer()
	for b.Loop() {
		mgr.Search("еСт")
	}
}

func BenchmarkStatsTasks(b *testing.B) {
	store := &MockStorage{tasks: []*task.Task{}}
	filter := &manager.FilterTasks{}
	render := &MockRender{}
	mgr := manager.NewManager(store, filter, render)
	testData := map[string]string{
		"title":       "Тест",
		"description": "Тест",
	}
	mgr.Create(testData)
	b.ResetTimer()
	for b.Loop() {
		mgr.Stats()
	}
}
