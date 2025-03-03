package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// URL уже существующего сервиса (например, запущенного на localhost:8080)
	targetURL, err := url.Parse("http://tasks-service:8080")
	if err != nil {
		log.Fatalf("Неверный URL: %v", err)
	}

	// Создаем обратный прокси
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	// При необходимости можно модифицировать Director, чтобы менять пути или заголовки запроса
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Например, можно добавить общие заголовки или логировать запросы
		log.Printf("Проксирование запроса: %s %s", req.Method, req.URL.String())
	}

	// Запускаем HTTP-сервер, который принимает запросы на порту 8000
	log.Println("API Gateway запущен на порту 8000...")
	if err := http.ListenAndServe(":8000", proxy); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
