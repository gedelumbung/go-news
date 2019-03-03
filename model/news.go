package model

import (
	"encoding/json"
	"time"

	"github.com/gedelumbung/go-news/helper"
	"github.com/go-sql-driver/mysql"
)

type News struct {
	ID        int            `db:"id"`
	Author    string         `db:"author"`
	Body      string         `db:"body"`
	CreatedAt mysql.NullTime `db:"created_at"`
}

type EsNews struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
}

func (o News) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        int    `json:"id"`
		Author    string `json:"author"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at"`
	}{
		ID:        o.ID,
		Author:    o.Author,
		Body:      o.Body,
		CreatedAt: helper.NullTimeToString(o.CreatedAt, time.RFC3339),
	})
}
