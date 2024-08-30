package file_manager

import (
	"fmt"
	"os"
	"sync"
)

// PrintFileDetails - Функция для формирования списка файлов и поддиректорий
func PrintFileDetails(mdir string) ([]FileEntry, error) {
	files, err := os.ReadDir(mdir)
	if err != nil { // Обработка ошибки
		return nil, fmt.Errorf("не удалось прочитать '%s'", err)
	}
	var sortedEntries []FileEntry = make([]FileEntry, len(files)) // Создание массива
	var wg sync.WaitGroup                                         //Инициализация WaitGroup
	for i, file := range files {                                  //Цикл для перебора файлов
		info, err := file.Info() // получения информации о файле
		if err != nil {          // Обработка ошибки
			fmt.Printf("Не удалось получить информацию о файле '%s': %v\n", file.Name(), err)
			continue
		}
		isDirectory := info.IsDir() //Проверка директории
		if isDirectory {            //Обработка директорий
			wg.Add(1) // Запускается новая горутина с функцией - GetDir
			go getDir(file, mdir, sortedEntries, i, &wg)
		} else {
			sortedEntries[i] = FileEntry{ // Обработка файлов
				FName: file.Name(),
				FSize: float64(info.Size()),
				FType: "Файл",
			}
		}
	}
	wg.Wait()                 //Ожидается завершение всех горутин
	return sortedEntries, nil //Возвращаются отсортированный массив и нил без ошибок
}

// GetDir - Функция для получения информации о поддиректории
func getDir(file os.DirEntry, pathDirectory string, filesEntry []FileEntry, index int, wg *sync.WaitGroup) {
	var size int64             // Инициализация переменной для хранения размера
	fileInfo, _ := file.Info() //Получение информации о файле
	defer wg.Done()
	// Рекурсивный вызов для поддиректорий
	directSum, err := calcSumDirect(fmt.Sprintf("%s/%s", pathDirectory, file.Name()))

	if err != nil { // Обработка ошибки
		fmt.Printf("Не удалось получить размер дирректории '%s': %v\n", file.Name(), err)
		return
	}
	size += directSum       // Добавление размера директории
	size += fileInfo.Size() // Добавление размера файла к общему размеру
	// Создание нового элемента FileEntry с информацией о директории
	filesEntry[index] = FileEntry{FName: file.Name(), FSize: float64(size), FType: "Директория"}
}

// calcSumDirect-Функция расчета суммы размеров файлов и поддиректорий
func calcSumDirect(pathDirectory string) (int64, error) {
	var sum int64                           // Инициализация переменной для накопления суммы размеров
	files, err := os.ReadDir(pathDirectory) // Чтение списка файлов и поддиректорий в директории
	if err != nil {                         // Обработка ошибки
		return 0, err
	}
	for _, file := range files {
		if file.IsDir() {
			// Рекурсивный вызов для поддиректорий
			dirSum, err := calcSumDirect(fmt.Sprintf("%s/%s", pathDirectory, file.Name()))
			if err != nil { //Обработка ошибки
				fmt.Printf("Не удалось получить размер директории '%s': %v\n", file.Name(), err)
				continue
			}
			sum += dirSum // Добавление размера директории
		} else {

			info, err := file.Info() // Получение информации о файле
			if err != nil {          //Обработка ошибки
				fmt.Printf("Не удалось получить размер директории '%s': %v\n", file.Name(), err)
				continue
			}
			sum += info.Size() // Размер добавляется к общей сумме
		}
	}
	return sum, nil
}
