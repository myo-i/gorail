package db

import (
	"database/sql"
	"fmt"
	"gorail/util"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Urls struct {
	Url       string
	Title     string
	LastVisit int
}

func GetData() {
	// 環境変数取得
	config, err := util.Load()
	if err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}

	// DB接続
	db, err := sql.Open("sqlite3", config.DbPath)
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}
	defer db.Close()

	// クエリを実行
	query := "SELECT title, url, last_visit_time FROM urls LIMIT 5"
	rows, err := db.Query(query)
	var data []Urls
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	// データ配列に格納
	for rows.Next() {
		var urls Urls
		err = rows.Scan(&urls.Url, &urls.Title, &urls.LastVisit)
		if err != nil {
			log.Fatalf("Failed to scan rows: %v", err)
		}
		data = append(data, urls)
	}
	for _, row := range data {
		fmt.Println(row)
	}

}
