package filemaneger

import (
	"sort"
)

// Структура для хранения информации о файле или директории
type FileEntry struct {
	FType   string
	FSize   float64
	FName   string
	TypSize string
}

// Функция форматирования массива FileEntry с использованием функции formatSize
func FormatFileEntries(dataFiles []FileEntry) {
	for i := range dataFiles {
		typ, razm := formatSize(dataFiles[i].FSize)
		dataFiles[i].TypSize = typ
		dataFiles[i].FSize = razm
	}
}

// Форматирование размера
func formatSize(bytes float64) (string, float64) {
	const b float64 = 1000
	// Сортируем по убыванию размеров
	switch {
	case bytes >= b*b*b: // Если размер в гигабайтах или больше
		return "Гигабайт", float64(bytes) / float64(b*b*b) // Форматирование до двух знаков после запятой
	case bytes >= b*b: // Если размер в мегабайтах
		return "Мегабайт", float64(bytes) / float64(b*b) // Форматирование до двух знаков после запятой
	case bytes >= b: // Если размер в килобайтах
		return "Килобайт", float64(bytes) / float64(b) // Форматирование до двух знаков после запятой
	default: // Если размер меньше килобайта
		return "Байт", float64(bytes) // Без форматирования для небольших размеров
	}
}

const Asc = "asc" // Создание констант
const Desc = "desc"

// Функция сортировки массива FileEntry по размеру файла
func SortFileEntry(dataFiles []FileEntry, SortFlag string) {
	// Определяем направление сортировки на основе параметра ask
	if SortFlag == Asc {
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
