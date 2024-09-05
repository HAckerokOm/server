package server

import (
	"bytes"
	"context"
	"encoding/json"
	"filespackage/serverr/file_manager"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// Send - Структура для пост запроса
type Send struct {
	RequestTime float64 `json:"f_requesttime"` // Время выполнения запроса
	Size        float64 `json:"f_size"` // Размер директории
	Root        string  `json:"f_root"` // Корневая директория
}

// handleFSRequest - Функция обработки запросов к файловой системе
func handleFSRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()// Время начала выполнения запроса
	// Структура для ответа 
	response := file_manager.Response{} 
	// Структура для отправки статистики на Apache сервер
	send := Send{} 

	//функция для отправки ответа клиенту и статистики на Apache сервер
	defer func() { 
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка при формировании JSON: %v", err), http.StatusInternalServerError)
			return
		}

		// Установка заголовка CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") 
		// Установка типа контента
		w.Header().Set("Content-Type", "application/json") 
		w.Write(jsonData) 

		// Время окончания выполнения запроса
		endTime := time.Now() 
		// Вычисление длительности выполнения
		duration := endTime.Sub(startTime) 
		// Запись времени в структуру Send
		send.RequestTime = duration.Seconds() 

		if response.Data != "" {
			// Отправка статистики на Apache сервер
			SendStatsToApache(send) 
		}
		fmt.Printf("Время выполнения запроса: %v\n", duration)
	}()
	// Получение корневой директории из переменных окружения
	defaultRoot := os.Getenv("DEFAULT_ROOT") 
	if defaultRoot == "" {
		// Обработка ошибки отсутствия корневой директории
		fmt.Println("Не указана директория") 
		return 
	}
	// Получение параметра dst из URL
	dst := r.URL.Query().Get("dst") 
	if dst == "" { 
		// Если dst не указан, устанавливаем значения по умолчанию
		response.Status = 200
		response.Error = ""
		response.Root = defaultRoot
		response.Data = " "
		return
	}
	// Получение параметра sort из URL
	sort := r.URL.Query().Get("sort") 

	// Получение списка файлов в директории
	entryfiles, err := file_manager.PrintFileDetails(dst) 
	// Обработка ошибки при получении списка файлов
	if err != nil { 
		response.Status = 500
		response.Error = fmt.Sprintf("ошибка при формировании списка файлов: %v", err)
	} else {
		// Запись пути в структуру Send
		send.Root = dst 
		// Переменная для суммирования размеров файлов
		var asum float64  
		for _, file := range entryfiles {
			// Суммирование размеров файлов
			asum += file.FSize 
		}
		// Запись общего размера в структуру Send
		send.Size = asum 

		 // Сортировка файлов
		file_manager.SortFileEntry(entryfiles, sort)
		// Форматирование данных о файлах
		file_manager.FormatFileEntries(entryfiles) 

		// Устанавливаем значения ответа
		response.Status = 200 
		response.Data = entryfiles
		response.Root = defaultRoot
	}

}

// StartServ - Функция для запуска сервера
func StartServ() {
	// Загрузка переменных окружения из .env файла
	err := godotenv.Load() 
	if err != nil {
		fmt.Println("Ошибка при загрузке переменных из .env файла")
	}
	// Получение порта из переменных окружения
	port := os.Getenv("PORT") 
	if port == "" {
		fmt.Println("Отсутствует обязательная переменная окружения PORT")
		os.Exit(1)
	}

	fs := http.FileServer(http.Dir("./resource")) 
	http.Handle("/", fs) //обработчик 
	http.HandleFunc("/fs", handleFSRequest)

	// Вывод информации о запуске сервера
	fmt.Printf("Сервер запущен на порту %s\n", port) 

	// Создание контекста с обработкой сигналов
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTSTP) 
	defer cancel() 
	// Адрес сервера
	httpServer := &http.Server{
		Addr: ":" + port, 
	}
	go func() {
		// Обработка ошибок при запуске сервера
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP-сервер прослушивает и обслуживает(): %v", err) 
		}
	}()
		// Ожидание завершения работы
	<-ctx.Done() 
	// Обработка ошибок при остановке сервера
	if err := httpServer.Shutdown(context.Background()); err != nil {
		fmt.Printf("Завершение работы HTTP-сервера(): %v", err) 
	}
}
// Функция отправки запроса на апаче
func SendStatsToApache(send Send) {
	// Получение ссылки из переменных окружения
	url := os.Getenv("URLAP") 
	if url == "" {
		fmt.Println("Отсутствует переменная окружения URLAP")
	}
	// Преобразование структуры в JSON
	jsonData, err := json.Marshal(send) 
	if err != nil {
		fmt.Printf("Ошибка при преобразовании данных в JSON: %v\n", err)
		return
	}

	// Отправка POST-запроса
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData)) 
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}
	// Закрытие тела ответа после использования
	defer resp.Body.Close() 

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body) 
	if err != nil {
		fmt.Printf("Ошибка при чтении тела ответа: %v\n", err)
		return
	}

	if len(body) == 0 {
		fmt.Println("Получен пустой ответ от сервера")
		return
	}

	var result map[string]interface{} // Мапа для парсинга JSON-ответа
	err = json.Unmarshal(body, &result) // Парсинг JSON
	if err != nil {
		fmt.Printf("Ошибка при парсинге JSON: %v\n", err)

		fmt.Printf("Необработанное тело ответа:\n%s\n", string(body))
		return
	}

	fmt.Printf("Ответ от Apache сервера: %+v\n", result) // Вывод полученного ответа
}
