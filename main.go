package main

import (
	"fmt"
	"gorail/config"
	"gorail/db"
	"gorail/user"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// 設定ファイル読み込み
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}

	// Chromeの履歴を読み込み
	srcDatabase, err := os.Open(getChromeHistoryPath())
	if err != nil {
		log.Fatalf("Failed to open source database: %v", err)
	}
	defer srcDatabase.Close()

	// データベースのコピー先の作成
	dstDatabase, err := os.Create(config.DbPath)
	if err != nil {
		log.Fatalf("Failed to create destination database: %v", err)
	}
	defer dstDatabase.Close()

	// コピー
	_, err = io.Copy(dstDatabase, srcDatabase)
	if err != nil {
		log.Fatalf("Failed to copy database: ", err)
	}
	fmt.Printf("Copied Chrome history succeeded!!")

	// DBからデータを取得
	data := db.GetData(config)

	// 滞在時間の長かったサイトのタイトルと滞在時間を取得
	topTenKey, topTenValue := db.GetLengthOfStay(data)

	// GUI起動
	user.RunApp(topTenKey, topTenValue)
}

func getChromeHistoryPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	historyPath := filepath.Join(homeDir, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "History")

	return historyPath
}
