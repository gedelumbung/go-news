package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/gedelumbung/go-news/helper"
	"github.com/gedelumbung/go-news/params"
	"github.com/gedelumbung/go-news/service"
	"github.com/gedelumbung/go-news/worker"
)

func NewsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getNews(w, r)
	case "POST":
		createNews(w, r)
	}
}

func getNews(w http.ResponseWriter, r *http.Request) {
	pageParams := r.URL.Query().Get("page")
	if pageParams == "" {
		pageParams = "1"
	}
	page := helper.StringToInt(pageParams)
	if page < 2 {
		page = 0
	} else {
		page = page - 1
	}
	start := page * 10

	data, err := service.GetNewsFromEs(start, 10)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err != nil {
		response, _ := json.Marshal(params.Response{
			err.Error(),
			nil,
			nil,
		})
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	sort.Slice(data, func(i int, j int) bool {
		return data[i].ID > data[j].ID
	})

	response, _ := json.Marshal(params.Response{
		"success",
		data,
		nil,
	})

	w.WriteHeader(200)
	w.Write(response)
	return
}

func createNews(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response, _ := json.Marshal(params.Response{
			err.Error(),
			nil,
			nil,
		})
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	var newsParams params.NewsRequest
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.Unmarshal(body, &newsParams)
	if err != nil {
		response, _ := json.Marshal(params.Response{
			err.Error(),
			nil,
			nil,
		})
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	worker.NewsWorker <- &params.NewsRequest{
		Author: newsParams.Author,
		Body:   newsParams.Body,
	}

	response, _ := json.Marshal(params.Response{
		"success",
		newsParams,
		nil,
	})

	w.WriteHeader(200)
	w.Write(response)
	return
}
