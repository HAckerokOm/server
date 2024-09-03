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
	RequestTime float64 `json:"f_requesttime"`
	Size        float64 `json:"f_size"`
	Root        string  `json:"f_root"`
}

// handleFSRequest - Функция для обработки запросов к файловой системе
func handleFSRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	response := file_manager.Response{}
	send := Send{}

	defer func() {
		jsonData, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, fmt.Sprintf("ошибка при формировании JSON: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)

		endTime := time.Now()
		duration := endTime.Sub(startTime)
		send.RequestTime = duration.Seconds()

		if response.Data != "" {
			SendStatsToApache(send)
		}
		fmt.Printf("Время выполнения запроса: %v\n", duration)
	}()

	defaultRoot := os.Getenv("DEFAULT_ROOT")
	if defaultRoot == "" {
		fmt.Println("Не указана директория")
		os.Exit(1)
	}

	dst := r.URL.Query().Get("dst")
	if dst == "" {
		response.Status = 200
		response.Error = ""
		response.Root = defaultRoot
		response.Data = " "
		return
	}

	sort := r.URL.Query().Get("sort")

	entryfiles, err := file_manager.PrintFileDetails(dst)
	if err != nil {
		response.Status = 500
		response.Error = fmt.Sprintf("ошибка при формировании списка файлов: %v", err)
	} else {
		send.Root = dst
		var asum float64
		for _, file := range entryfiles {
			asum += file.FSize
		}
		send.Size = asum

		file_manager.SortFileEntry(entryfiles, sort)
		file_manager.FormatFileEntries(entryfiles)

		response.Status = 200
		response.Data = entryfiles
		response.Root = defaultRoot
	}

}

// StartServ - Функция для запуска сервера
func StartServ() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка при загрузке переменных из .env файла")
	}

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("Отсутствует обязательная переменная окружения PORT")
	}

	fs := http.FileServer(http.Dir("./resource"))
	http.Handle("/", fs)
	http.HandleFunc("/fs", handleFSRequest)

	fmt.Printf("Сервер запущен на порту %s\n", port)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTSTP)
	defer cancel()

	httpServer := &http.Server{
		Addr: ":" + port,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server ListenAndServe(): %v", err)
		}
	}()

	<-ctx.Done()

	if err := httpServer.Shutdown(context.Background()); err != nil {
		fmt.Printf("HTTP server Shutdown(): %v", err)
	}
}

func SendStatsToApache(send Send) {
	url := "http://localhost:80/post.php"

	jsonData, err := json.Marshal(send)
	if err != nil {
		fmt.Printf("Ошибка при преобразовании данных в JSON: %v\n", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении тела ответа: %v\n", err)
		return
	}

	if len(body) == 0 {
		fmt.Println("Получен пустой ответ от сервера")
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Ошибка при парсинге JSON: %v\n", err)

		fmt.Printf("Raw response body:\n%s\n", string(body))
		return
	}

	fmt.Printf("Ответ от Apache сервера: %+v\n", result)
}
