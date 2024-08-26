package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)
// Функция расчета суммы размеров файлов и поддиректорий
func calcSumDirect(pathDirectory string) (int64, error) {
	var sum int64 // Инициализация переменной для накопления суммы размеров
	files, err := os.ReadDir(pathDirectory) // Чтение списка файлов и поддиректорий в директории
	if err != nil {
		return 0, err
	}
	for _, file := range files {
		if file.IsDir() {
			// Рекурсивный вызов для поддиректорий
			dirSum, err := calcSumDirect(fmt.Sprintf("%s/%s", pathDirectory, file.Name()))
			if err != nil {
				fmt.Printf("Не удалось получить размер директории '%s': %v\n", file.Name(), err)
				continue
			}
			sum += dirSum
		} else {
			// Получение информации о файле
			info, _ := file.Info()
			sum += info.Size()
		}
	}
	return sum, nil
}
// Функция для получения информации о поддиректории
func GetDir(file os.DirEntry, pathDirectory string, filesEntry []FileEntry, index int, wg *sync.WaitGroup) {
	var size int64 // Инициализация переменной для хранения размера
	fileInfo, _ := file.Info()
	defer wg.Done()
	// Рекурсивный вызов для поддиректорий0
	directsum, err := calcSumDirect(fmt.Sprintf("%s/%s", pathDirectory, file.Name()))
	size += directsum
	if err != nil {
		fmt.Printf("Не удалось получить размер дирректории '%s': %v\n", file.Name(), err)
		return
	}
	size += fileInfo.Size() // Добавление размера файла к общему размеру
	// Создание нового элемента FileEntry с информацией о директории
	filesEntry[index] = FileEntry{FName: file.Name(), FSize: size, FType: "Директория"}
}
// Функция для формирования списка файлов и поддиректорий
func printFileDetails(mdir string) []FileEntry {
	files, _ := os.ReadDir(mdir)
	var sortedEntries []FileEntry = make([]FileEntry, len(files))
	var wg sync.WaitGroup
	fmt.Printf("%-15s %-15s %-30s\n", "Тип", "Размер", "Название") // Вывод заголовков столбцов
	for i, file := range files {
		info, err := file.Info()
		if err != nil {
			fmt.Printf("Не удалось получить информацию о файле '%s': %v\n", file.Name(), err)
			continue
		}
		isDirectory := info.IsDir()
		if isDirectory {
			wg.Add(1)
			go GetDir(file, mdir, sortedEntries, i, &wg)
		} else {
			sortedEntries[i] = FileEntry{
				FName: file.Name(),
				FSize: info.Size(),
				FType: "Файл",
			}
		}
	}
	wg.Wait()
	return sortedEntries
}

func main() {
	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка при загрузке переменных из .env файла")
	}

	// Читаем значения переменных
	port := os.Getenv("PORT")
	if port == "" {
        fmt.Println("Отсутствует обязательная переменная окружения PORT")
    }

	fmt.Printf("Сервер запускается на порту %s...\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dst := r.URL.Query().Get("dst")
		sort := r.URL.Query().Get("sort")
		var sortAscending bool
		if sort == "asc" {
			sortAscending = true
		} else if sort == "desc" {
			sortAscending = false
		} else {
			sortAscending = true
			fmt.Printf("Не выбран флаг сортировки, будет использован asc \n")
		}

		entryfiles := printFileDetails(dst)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("ошибка при формировании списка файлов: %v", err)))
			return
		}
		SortFileEntry(entryfiles, sortAscending)

		jsonData, err := json.MarshalIndent(entryfiles, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("ошибка при форматировании JSON: %v", err)))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})
	// Запуск сервера
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("ошибка запуска сервера: %v", err)
	}

	startTime := time.Now()

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Время выполнения программы: %v\n", duration)
}
