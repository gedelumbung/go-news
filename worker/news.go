package worker

import (
	"fmt"
	"time"

	"github.com/gedelumbung/go-news/model"
	"github.com/gedelumbung/go-news/params"
	"github.com/gedelumbung/go-news/service"
)

func CreateNewsWorker(jobs <-chan *params.NewsRequest) {
	fmt.Println("Register the worker")
	for i := range jobs {
		data := model.News{
			Author: i.Author,
			Body:   i.Body,
		}
		newsDb, _ := service.StoreNewsToDb(&data)
		newsEs, _ := service.StoreNewsToEs(newsDb)
		fmt.Println("worker processing job (MySQL) with #ID : ", newsDb.ID)
		fmt.Println("worker processing job (ES) with #ID : ", newsEs.ID)
		time.Sleep(time.Second * 5)
	}
}
