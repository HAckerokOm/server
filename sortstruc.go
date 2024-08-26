package main

import (
	"fmt"
	"sort"
)
// Структура для хранения информации о файле или директории
type FileEntry struct {
	FType string
	FSize int64
	FName string
}

// Функция для печати информации о файле или директории
func (fileEntry *FileEntry) Print() {
	// Получаем форматированный размер файла
	newSize := formatSize(fileEntry.FSize)
	// Выводим информацию в формате:
	// Тип     Размер      Название
	fmt.Printf("%-15s %-16s %-10s\n", fileEntry.FType, newSize, fileEntry.FName)
}
// Форматирование размера
func formatSize(bytes int64) string {
	// Сортируем по убыванию размеров
	switch {
	case bytes >= 1000*1000*1000: // Если размер в гигабайтах или больше
		return fmt.Sprintf("%.2f Гигабайт", float64(bytes)/float64(1000*1000*1000))// Форматирование до двух знаков после запятой
	case bytes >= 1000*1000: // Если размер в мегабайтах
		return fmt.Sprintf("%.2f Мегабайт", float64(bytes)/float64(1000*1000))// Форматирование до двух знаков после запятой
	case bytes >= 1000: // Если размер в килобайтах
		return fmt.Sprintf("%.2f Килобайт", float64(bytes)/float64(1000))// Форматирование до двух знаков после запятой
	default: // Если размер меньше килобайта
		return fmt.Sprintf("%d Байт", bytes)// Без форматирования для небольших размеров
	}
}
// Функция сортировки массива FileEntry по размеру файла
func SortFileEntry(dataFiles []FileEntry, ask bool) {
	// Определяем направление сортировки на основе параметра ask
	if !ask {
		// Сортировка по возрастанию размера
		sort.Slice(dataFiles, func(i, j int) bool {
			return dataFiles[i].FSize < dataFiles[j].FSize
		})
	} else {
		// Сортировка по убыванию размера
		sort.Slice(dataFiles, func(i, j int) bool {
			return dataFiles[i].FSize > dataFiles[j].FSize
		})
	}
}
