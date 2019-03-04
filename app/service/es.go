package service

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/gedelumbung/go-news/model"
	"github.com/olivere/elastic"
)

const EsIndex = "news_index"
const EsType = "news"
const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"news":{
			"properties":{
				"created_at":{
					"type":"date",
					"format":"YYYY-MM-DD'T'HH:mm:ss"
				}
			}
		}
	}
}`

var esNewsQueue chan *model.EsNews

func connectEs() *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL(os.Getenv("ES_URL")))
	if err != nil {
		panic("Could not connect to the elasticsearch")
	}
	//check index if exist
	exists, _ := client.IndexExists(EsIndex).Do(context.Background())

	//create index if not exist
	if !exists {
		_, err := client.CreateIndex(EsIndex).BodyString(mapping).Do(context.Background())
		if err != nil {
			panic("Could not create index")
		}
	}
	return client
}

func StoreNewsToEs(data *model.News) (*model.News, error) {
	client := connectEs()

	t := time.Now()
	timeFormatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	newsParams := &model.EsNews{ID: data.ID, CreatedAt: timeFormatted}
	_, err := client.Index().
		Index(EsIndex).
		Type(EsType).
		Id(strconv.Itoa(data.ID)).
		BodyJson(newsParams).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	return data, err
}

func GetNewsFromEs(start, limit int) ([]*model.News, error) {
	client := connectEs()
	data, err := client.Search().
		Index(EsIndex).
		Sort("created_at", false).
		From(start).Size(limit).
		Pretty(true).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	var (
		esNewsItem model.EsNews
		news       []*model.News
		wg         sync.WaitGroup
	)
	news = []*model.News{}
	esNewsQueue = make(chan *model.EsNews, 100)

	for _, item := range data.Each(reflect.TypeOf(esNewsItem)) {
		if e, ok := item.(model.EsNews); ok {
			wg.Add(1)
			go func(e model.EsNews) {
				esNewsQueue <- &model.EsNews{
					ID:        e.ID,
					CreatedAt: e.CreatedAt,
				}
			}(e)
		}
	}

	go func() {
		for e := range esNewsQueue {
			data, _ := FindNewsByID(e.ID)
			news = append(news, &data)
			wg.Done()
		}
	}()

	wg.Wait()

	return news, err
}
