package server

import (
	"context"
	"encoding/json"
	"filespackage/file_manager"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const DEFAULT_ROOT string = "/"

// handleFSRequest-Функция обработки  http запроса
func handleFSRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now() // Записываем текущее время для измерения времени выполнения запроса
	response := file_manager.Response{}
	defer func() {
		jsonData, err := json.MarshalIndent(response, "", "  ") // Переводим структуру response в JSON с отступами
		if err != nil {                                         // Обработка ошибки
			http.Error(w, fmt.Sprintf("ошибка при форматировании JSON: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")     // Устанавливаем заголовок Access-Control-Allow-Origin для CORS
		w.Header().Set("Content-Type", "application/json")     // Устанавливаем заголовок типа содержимого ответа
		w.Write(jsonData)                                      // Отправляем данные клиенту
		endTime := time.Now()                                  // Получаем текущее время
		duration := endTime.Sub(startTime)                     // Вычисляем разницу между конечным и начальным временем
		fmt.Printf("Время выполнения запроса: %v\n", duration) // Время выполнения запроса
	}()
	dst := r.URL.Query().Get("dst")   // Получаем путь назначения из URL-запроса
	sort := r.URL.Query().Get("sort") // Получаем флаг сортировки из URL-запроса

	defaultRoot := os.Getenv("DEFAULT_ROOT")
	if defaultRoot == "" {
		defaultRoot = DEFAULT_ROOT
	}

	entryfiles, err := file_manager.PrintFileDetails(dst) // Вызываем функцию для получения  файлов и директорий
	if err != nil {                                       // Проверяем наличие ошибок при получении файловых записей
		response.Status = 500 // Если ошибка существует, создаем структуру Response с кодом статуса 500 (Internal Server Error)
		response.Error = fmt.Sprintf("ошибка при формировании списка файлов: %v", err)
		response.Data = ""
	}

	file_manager.SortFileEntry(entryfiles, sort) // Вызываем функцию сортировки файловых записей
	file_manager.FormatFileEntries(entryfiles)   // Вызываем функцию форматирования размеров файловых записей

	response.Status = 200      // Status-устанавливается в 200 для успешного ответа
	response.Error = ""        // Error - оставляется пустым, так как нет ошибок
	response.Data = entryfiles // Data-содержит список файловых записей entryfiles
}

// StartServ-функция запуска сервера
func StartServ() {

	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil { // Обработка ошибки
		fmt.Println("Ошибка при загрузке переменных из .env файла")
	}

	// Читаем значения переменных
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Отсутствует обязательная переменная окружения PORT")
	}

	defaultRoot := os.Getenv("DEFAULT_ROOT")
	if defaultRoot == "" {
		defaultRoot = DEFAULT_ROOT
	}

	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./static/style"))))
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("./static/script"))))
	http.HandleFunc("/", serveHTMLPage)
	http.HandleFunc("/fs", handleFSRequest) //обработчик http-запросов
	fmt.Printf("Сервер запущен на порту %s\n", port)
	// Закрытие сервера

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTSTP) //контекст для отслеживания сигналов
	defer cancel()
	// Здесь логика сервера
	httpServer := &http.Server{
		Addr: ":" + port,
	}
	go func() { // запуск в отдельной горутине
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed { // Обработка ошибки
			fmt.Printf("HTTP server ListenAndServe(): %v", err)
		}
	}()
	// Ожидание завершения работы сервера
	<-ctx.Done()

	// Корректное завершение работы HTTP сервера
	if err := httpServer.Shutdown(context.Background()); err != nil { // Обработка ошибки
		fmt.Printf("HTTP server Shutdown(): %v", err)
	}
}

// serveHTMLPage-функция для обработки HTTP-запросов на HTML страницу
func serveHTMLPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/site.html")

}
