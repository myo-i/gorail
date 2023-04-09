package db

import (
	"database/sql"
	"fmt"
	"gorail/config"
	"log"
	"regexp"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DBからデータを取得した際の構造体
type TimeOnSite struct {
	Title         string
	Url           string
	VisitDuration int
}

// とりあえず3か月分くらいのデータを（できれば滞在時間が長い順で）取得する
// 取得したデータを1日、1週間、1か月とかでさらに範囲を絞って使う

type History interface {
	GetHistory()
}

// DBからデータを取得
func GetData() []TimeOnSite {
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
	// query := "SELECT title, url, last_visit_time FROM urls LIMIT 5"
	// 1ヶ月(30日)は2592000秒
	query := "SELECT urls.title, urls.url, visits.visit_duration FROM visits LEFT JOIN urls on visits.url = urls.id WHERE urls.last_visit_time >= (strftime('%s', 'now', '-5 days')+11644473600)*1000000 ORDER BY visits.id desc LIMIT 10;"

	// query := "SELECT COUNT(title) FROM urls"
	rows, err := db.Query(query)
	var data []TimeOnSite
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	// データ配列に格納
	for rows.Next() {
		var urls TimeOnSite
		err = rows.Scan(&urls.Title, &urls.Url, &urls.VisitDuration)
		if err != nil {
			log.Fatalf("Failed to scan rows: %v", err)
		}
		data = append(data, urls)
	}
	for _, row := range data {
		fmt.Println(row)
	}
	return data

}

// goroutineを使って並列で処理
// ホスト名をキーに、visit_durationを足していく（できればソートも）
func CalcTimeOnSite(datas []TimeOnSite) sync.Map {
	// スライスとして宣言するかは要検討
	// var hostAndTime = make(map[string]int)
	var hostAndTime = sync.Map{}
	var wg sync.WaitGroup
	// c := make(chan bool)

	// データの個数分goroutineを実行するので、Addにはdatasの要素数を設定
	// wg.Done()が実行されるとAddが減っていく
	// var clientMutex sync.Mutex

	wg.Add(10)
	for _, val := range datas {
		go createTimeOnSiteMap(val, hostAndTime, &wg)
	}

	// Addが0になるまで待つ
	wg.Wait()
	normal("hello")

	// 何故かチャネル追加したらhostAndTimeに値が入ったからコード読んでおく
	// fmt.Println(hostAndTime)
	// <-c
	// fmt.Println(hostAndTime)

	return hostAndTime
}

// mutexはコストが高く、十分なスケーラビリティを確保することができないため、Go1.9以降のバージョンでは、sync.Mapという並行安全なマップが追加されました[2]。
//
//	func createTimeOnSiteMap(val TimeOnSite, hostAndTime sync.Map, mutex sync.Mutex, wg *sync.WaitGroup) {
//		mutex.Lock()
//		hostAndTime[urlToHostName(val.Url)] += val.VisitDuration
//		defer mutex.Unlock()
//		fmt.Println(val.Title)
//		defer wg.Done()
//		// fmt.Println(hostAndTime)
//		// 何故か↓のチャネルを追記したらhostAndTimeに値が入った
//		// c <- true
//	}
func createTimeOnSiteMap(val TimeOnSite, hostAndTime sync.Map, wg *sync.WaitGroup) {
	hostAndTime.Store(val.Url, val.VisitDuration)
	fmt.Println(val.Title)
	defer wg.Done()
	// fmt.Println(hostAndTime)
	// 何故か↓のチャネルを追記したらhostAndTimeに値が入った
	// c <- true
}

// 取得した際のタイトルはバラバラなのでURLで判断
// ホスト名までをキーバリューに格納してvisit_durationを足していく、そして合計を分や時間に直す
// URLは(http|https)://[\w\.-]+/でホスト名までにしてキーにする
func urlToHostName(url string) string {
	rex := regexp.MustCompile("(http|https)://[\\w\\.-]+/")
	return rex.FindString(url)
}

func normal(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}
