package service

import (
	"database/sql"
	"errors"
	"os"

	"github.com/gedelumbung/go-news/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound = errors.New("item not found")
)

func connectDb() *sqlx.DB {
	db, err := sqlx.Open("mysql", os.Getenv("MYSQL_URL"))
	db.SetMaxOpenConns(100)
	if err != nil {
		panic("Could not connect to the db")
	} else {
		return db
	}
}

func StoreNewsToDb(model *model.News) (*model.News, error) {
	db := connectDb()
	stmt, err := db.Preparex(`insert into news
		(author, body)
		values (?, ?)`)

	defer db.Close()

	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(
		model.Author,
		model.Body,
	)

	id, err := result.LastInsertId()
	model.ID = int(id)

	if err != nil {
		return model, err
	}

	return model, err
}

func FindNewsByID(id int) (model.News, error) {
	var news model.News
	db := connectDb()

	err := db.QueryRowx(`select id, author, body, created_at from news where news.id = ?`, id).StructScan(&news)

	defer db.Close()

	if err == sql.ErrNoRows {
		return news, ErrNotFound
	}
	if err != nil {
		return news, err
	}
	return news, err
}
