package storage

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrSerializeJson   = errors.New("сериализации данных в json")
	ErrDeserializeJson = errors.New("преобразование данных в структуру")
)

// Дженерик-функция принимающий на вход любую структуру
func DataToJson[T any](data *T) ([]byte, error) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("ошибка: %w. Данные: %v", ErrSerializeJson, data)
	}
	return jsonData, nil
}

// указатель на jsonData не нужен - слайс уже в себе содержит указатель на массив - копируется только 24 байта
func JsonToData[T any](jsonData []byte, data *T) error {
	err := json.Unmarshal(jsonData, data)
	if err != nil {
		return fmt.Errorf("ошибка: %w. Данные: %v", ErrSerializeJson, data)
	}
	return nil
}
