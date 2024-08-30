package server

import (
	"context"
	"encoding/json"
	"filespackage/serverr/file_manager"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// handleFSRequest - Функция для обработки запросов к файловой системе
func handleFSRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()             // Запускаем таймер для измерения времени выполнения запроса
	response := file_manager.Response{} // Создаем пустой объект для хранения ответа

	defer func() { // Функция для очистки ресурсов после завершения запроса
		jsonData, err := json.MarshalIndent(response, "", "  ") // Преобразование ответа в JSON с отступами
		if err != nil {                                         // В случае ошибки при форматировании JSON отправляем ошибку пользователям
			http.Error(w, fmt.Sprintf("ошибка при форматировании JSON: %v", err), http.StatusInternalServerError)
			return
		}
		// Устанавливаем заголовки ответа
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData) // Отправляем данные пользователям
		endTime := time.Now()
		duration := endTime.Sub(startTime)                     // Рассчитываем время выполнения запроса
		fmt.Printf("Время выполнения запроса: %v\n", duration) // Логируем время выполнения запроса
	}()

	defaultRoot := os.Getenv("DEFAULT_ROOT") // Получаем значение переменной окружения DEFAULT_ROOT
	if defaultRoot == "" {                   // Выводим сообщение об ошибке, если переменная окружения не установлена
		fmt.Printf("Не указана дирректория")
		os.Exit(1)
	}
	// Получаем параметры запроса
	dst := r.URL.Query().Get("dst")
	sort := r.URL.Query().Get("sort")

	// Получаем список файлов и директорий
	entryfiles, err := file_manager.PrintFileDetails(dst)
	if err != nil {
		response.Status = 500 // Устанавливаем статус ошибки и сообщение об ошибке
		response.Error = fmt.Sprintf("ошибка при формировании списка файлов: %v", err)
	} else {
		// Сортируем и форматируем список файлов
		file_manager.SortFileEntry(entryfiles, sort)
		file_manager.FormatFileEntries(entryfiles)
		// Устанавливаем статус успешного выполнения
		response.Status = 200
		// Добавляем данные и корневую директорию в ответ
		response.Data = entryfiles
		response.Root = defaultRoot
	}
}

// serveHTMLPage - Функция для обслуживания HTML страниц
func serveHTMLPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/site.html") // Отправляем HTML страницу
}

// StartServ - Функция для запуска сервера
func StartServ() {
	err := godotenv.Load() // Загружаем переменные окружения из .env файла
	if err != nil {
		fmt.Println("Ошибка при загрузке переменных из .env файла")
	}

	port := os.Getenv("PORT") // Получаем порт из переменной окружения
	if port == "" {
		fmt.Println("Отсутствует обязательная переменная окружения PORT")
	}

	// Настроиваем маршруты для  файлов
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./static/style"))))
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("./static/script"))))

	// Определяем обработчики HTTP запросов
	http.HandleFunc("/", serveHTMLPage)
	http.HandleFunc("/fs", handleFSRequest)

	// Выводим сообщение о запуске сервера
	fmt.Printf("Сервер запущен на порту %s\n", port)

	// Создаем контекст для отслеживания сигналов
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTSTP)
	defer cancel()

	// Создаем объект сервера
	httpServer := &http.Server{
		Addr: ":" + port,
	}
	go func() { // Запускаем сервер в отдельном горутине
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server ListenAndServe(): %v", err)
		}
	}()

	// Ждем сигнал завершения
	<-ctx.Done()

	// Завершаем сервер
	if err := httpServer.Shutdown(context.Background()); err != nil {
		fmt.Printf("HTTP server Shutdown(): %v", err)
	}
}
