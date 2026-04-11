package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SleepData - та самая структура для хранения настроек
type SleepData struct {
	Hour   string `json:"hour"`
	Minute string `json:"minute"`
	Cycles string `json:"cycles"`
}

// getConfigPath определяет, где будет лежать файл сохранения.
// os.UserConfigDir() автоматически найдет папку AppData/Roaming на Windows.
func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return  "", err
	}
	// Создаем папку NekoSleep в AppData, если её нет
	appDir := filepath.Join(configDir, "NekoSleep")
	os.MkdirAll(appDir, os.ModePerm)
	
	return filepath.Join(appDir, "config.json"), nil
}

// Save сохраняет структуру в файл
func Save(data *SleepData) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	// Превращаем структуру в красивый JSON с отступами
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Записываем в файл (0644 - стандартные права доступа)
	return os.WriteFile(path, fileData, 0644)
}

// Load читает данные из файла
func Load() (*SleepData, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err // Файла еще нет или ошибка чтения
	}

	var data SleepData
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// Delete удаляет файл конфигурации
func Delete() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "NekoSleep", "config.json")
	return os.Remove(path)
}