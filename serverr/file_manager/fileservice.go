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
	// Создание массива
	var sortedEntries []FileEntry = make([]FileEntry, len(files)) 
	//Инициализация WaitGroup
	var wg sync.WaitGroup   
	//Цикл для перебора файлов                                      
	for i, file := range files {  
		// получения информации о файле                                
		info, err := file.Info() 
		// Обработка ошибки
		if err != nil {          
			fmt.Printf("Не удалось получить информацию о файле '%s': %v\n", file.Name(), err)
			continue
		}
		//Проверка директории
		isDirectory := info.IsDir() 
		 //Обработка директорий
		if isDirectory {    
			// Запускается новая горутина с функцией - GetDir       
			wg.Add(1) 
			go getDir(file, mdir, sortedEntries, i, &wg)
		} else {
			// Обработка файлов
			sortedEntries[i] = FileEntry{ 
				FName: file.Name(),
				FSize: float64(info.Size()),
				FType: "Файл",
			}
		}
	}
	//Ожидается завершение всех горутин
	wg.Wait()   
	//Возвращаются отсортированный массив и нил без ошибок              
	return sortedEntries, nil 
}

// GetDir - Функция для получения информации о поддиректории
func getDir(file os.DirEntry, pathDirectory string, filesEntry []FileEntry, index int, wg *sync.WaitGroup) {
	// Инициализация переменной для хранения размера
	var size int64    
	//Получение информации о файле         
	fileInfo, _ := file.Info() 
	defer wg.Done()
	// Рекурсивный вызов для поддиректорий
	directSum, err := calcSumDirect(fmt.Sprintf("%s/%s", pathDirectory, file.Name()))
	// Обработка ошибки
	if err != nil { 
		fmt.Printf("Не удалось получить размер дирректории '%s': %v\n", file.Name(), err)
		return
	}
	 // Добавление размера директории
	size += directSum  
	// Добавление размера файла к общему размеру    
	size += fileInfo.Size() 
	// Создание нового элемента FileEntry с информацией о директории
	filesEntry[index] = FileEntry{FName: file.Name(), FSize: float64(size), FType: "Директория"}
}

// calcSumDirect-Функция расчета суммы размеров файлов и поддиректорий
func calcSumDirect(pathDirectory string) (int64, error) {
	// Инициализация переменной для накопления суммы размеров
	var sum int64         
	// Чтение списка файлов и поддиректорий в директории                  
	files, err := os.ReadDir(pathDirectory) 
	// Обработка ошибки
	if err != nil {                         
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
			// Получение информации о файле
			info, err := file.Info() 
			 //Обработка ошибки
			if err != nil {         
				fmt.Printf("Не удалось получить размер директории '%s': %v\n", file.Name(), err)
				continue
			}
			// Размер добавляется к общей сумме
			sum += info.Size() 
		}
	}
	return sum, nil
}
