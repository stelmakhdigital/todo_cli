package storage

import (
	"os"
)

const fileName = "tasks.json"

const (
	permWrite = 2 << iota // 2
	permRead              // 4
)

const (
	owner = (permRead | permWrite) << 6 // 0600
	group = permRead << 3               // 0040
	other = permRead                    // 0004
)

const (
	fileMode644 = owner | group | other // 0644
)

type Storage interface {
	Save(tasks []byte) error
	// Load() ([]byte, error)
}

func Save(tasks []byte) error {
	return os.WriteFile(fileName, tasks, fileMode644)
}
