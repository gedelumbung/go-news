package main

import (
	"net/http"

	"github.com/joho/godotenv"

	"github.com/gedelumbung/go-news/handler"
	"github.com/gedelumbung/go-news/params"
	"github.com/gedelumbung/go-news/worker"
)

func main() {
	loadConfig()
	worker.NewsWorker = make(chan *params.NewsRequest, 100)
	go worker.CreateNewsWorker(worker.NewsWorker)

	http.HandleFunc("/news", handler.NewsHandler)
	http.ListenAndServe(":9090", nil)
}

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}
