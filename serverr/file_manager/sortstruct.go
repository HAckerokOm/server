package file_manager

import (
	"sort"
)
// Создание констант
const asc = "asc"   // конст сортировки по возрастанию
const desc = "desc" // конст сортировки по убыванию
const b float64 = 1000 // конст для форматир размера

// FileEntry - Структура для хранения информации о файле или директории
type FileEntry struct {
	FType    string  `json:"f_type"`     // Тип файл/директ
	FSize    float64 `json:"f_size"`     // размер
	FName    string  `json:"f_name"`     // название
	TypeSize string  `json:"f_typesize"` // тип размера(мб/кб/б/гб)
}

// Response - Структура ответов для обработки ошибок
type Response struct {
	Status int    `json:"f_status"` // Статус ответа
	Error  string `json:"f_error"`  // ошибка
	Data   any    `json:"f_data"`   // массив данных
	Root   string `json:"f_root"`
}

// FormatFileEntries - Функция форматирования массива FileEntry с использованием функции formatSize
func FormatFileEntries(dataFiles []FileEntry) {
	for i := range dataFiles {
		types, size := formatSize(dataFiles[i].FSize)
		dataFiles[i].TypeSize = types
		dataFiles[i].FSize = size
	}
}

// formatSize - Функция Форматирование размера
func formatSize(bytes float64) (string, float64) {
	
	// Сортируем по убыванию размеров
	switch {
	case bytes >= b*b*b: // Если размер в гигабайтах или больше
		return "Гигабайт", float64(bytes) / (b * b * b) // Форматирование до двух знаков после запятой
	case bytes >= b*b: // Если размер в мегабайтах
		return "Мегабайт", float64(bytes) / (b * b) // Форматирование до двух знаков после запятой
	case bytes >= b: // Если размер в килобайтах
		return "Килобайт", float64(bytes) / (b) // Форматирование до двух знаков после запятой
	default: // Если размер меньше килобайта
		return "Байт", float64(bytes) // Без форматирования для небольших размеров
	}
}



// SortFileEntry - Функция сортировки массива FileEntry по размеру файла
func SortFileEntry(dataFiles []FileEntry, sortFlag string) {
	// Определяем направление сортировки на основе параметра ask
	if sortFlag == asc {
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
