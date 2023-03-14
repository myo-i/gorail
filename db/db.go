package db

import (
	"database/sql"
	"fmt"
	"gorail/config"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Urls struct {
	Url       string
	Title     string
	LastVisit int
}

// とりあえず3か月分くらいのデータを（できれば滞在時間が長い順で）取得する
// 取得したデータを1日、1週間、1か月とかでさらに範囲を絞って使う

type History interface {
	GetHistory()
}

// DBからデータを取得
func GetData() {
	// 環境変数取得
	config, err := config.Load()
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
