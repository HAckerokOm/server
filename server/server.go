package server

import (
	"context"
	"encoding/json"
	"filespackage/filemaneger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func handleFSRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now() // Записываем текущее время для измерения времени выполнения запроса
	dst := r.URL.Query().Get("dst") // Получаем путь назначения из URL-запроса
	sort := r.URL.Query().Get("sort") // Получаем флаг сортировки из URL-запроса

	entryfiles, err := filemaneger.PrintFileDetails(dst)// Вызываем функцию для получения  файлов и директорий
	if err != nil { //Обработка ошибки
		http.Error(w, fmt.Sprintf("ошибка при формировании списка файлов: %v", err), http.StatusInternalServerError)
		return
	}
	filemaneger.SortFileEntry(entryfiles, sort) // Вызываем функцию сортировки файловых записей
	filemaneger.FormatFileEntries(entryfiles) // Вызываем функцию форматирования размеров файловых записей
	jsonData, err := json.MarshalIndent(entryfiles, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка при форматировании JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Printf("Время выполнения запроса: %v\n", duration)
}
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

	http.HandleFunc("/fs", handleFSRequest)
	fmt.Printf("Сервер запущен на порту %s\n", port)
	// Закрытие сервера
	Closer := make(chan os.Signal, 1) // Буфер размером 1, чтобы отправитель не блокировался
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTSTP)

	go func() {
		<-Closer // Ожидание сигнала
		cancel() // Отмена контекста
	}()
	// Здесь логика сервера
	httpServer := &http.Server{
		Addr: ":" + port,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed { // Обработка ошибки
			fmt.Printf("HTTP server ListenAndServe(): %v", err)
		}
	}()
	// Ожидание завершения работы сервера
	<-ctx.Done()
	// Закрытие канала Closer
	close(Closer)
	// Корректное завершение работы HTTP сервера
	if err := httpServer.Shutdown(context.Background()); err != nil { // Обработка ошибки
		fmt.Printf("HTTP server Shutdown(): %v", err)
	}
}
