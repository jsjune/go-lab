package main

import (
	"go-sample-stdlib-http/handler"
	"go-sample-stdlib-http/store"
	"log"
	"net/http"
	"time"
)

// loggingMiddleware — Gin의 Logger()와 동일한 역할, 직접 구현
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	s := store.New()
	itemHandler := handler.NewItemHandler(s)

	mux := http.NewServeMux()

	// health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	itemHandler.RegisterRoutes(mux)

	// 미들웨어는 핸들러를 감싸는 방식으로 적용
	// Gin의 r.Use()와 동일한 개념, 단지 명시적으로 작성
	server := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}

	log.Println("stdlib HTTP 서버 시작 — :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
