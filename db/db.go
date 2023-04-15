package db

import (
	"database/sql"
	"fmt"
	"gorail/config"
	"log"
	"regexp"
	"sync"

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
	query := "SELECT urls.title, urls.url, visits.visit_duration FROM visits LEFT JOIN urls on visits.url = urls.id WHERE urls.last_visit_time >= (strftime('%s', 'now', '-1 months')+11644473600)*1000000 ORDER BY visits.id desc;"

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
	// for _, row := range data {
	// 	fmt.Println(row)
	// }
	return data

}

func CalcTimeOnSite(datas []TimeOnSite) sync.Map {
	var hostAndTime = sync.Map{}
	var wg sync.WaitGroup

	// データの個数分goroutineを実行するので、Addにはdatasの要素数を設定
	// wg.Done()が実行されるとAddが減っていく
	// var clientMutex sync.Mutex
	wg.Add(len(datas))
	for _, data := range datas {
		// バリューが上書きされてしまう
		go func(data TimeOnSite) {
			value, ok := hostAndTime.Load(urlToHostName(data.Url))
			// バリューが上書きされないように処理
			if ok {
				hostAndTime.Store(urlToHostName(data.Url), data.VisitDuration+value.(int))
			} else {
				hostAndTime.Store(urlToHostName(data.Url), data.VisitDuration)
			}
			// fmt.Println(urlToHostName(data.Url), data.VisitDuration)
			defer wg.Done()
		}(data)
	}

	// Addが0になるまで待つ
	wg.Wait()

	// fmt.Println("Hello")
	// hostAndTime.Range(func(key interface{}, value interface{}) bool {
	// 	fmt.Printf("Key: %v(Type: %T) -> Value: %v(Type: %T)\n", key, key, value, value)
	// 	return true
	// })

	searchTopFive(hostAndTime)

	// mapをソートしたい
	// sortedMap := make(map[string]int)
	// keys := make([]string, 0, len(sortedMap))
	// hostAndTime.Range(func(key interface{}, value interface{}) bool {
	// 	sortedMap[key.(string)] = value.(int)
	// 	keys = append(keys, key.(string))
	// 	return true
	// })

	// sort.Slice(keys, func(i, j int) bool {
	// 	fmt.Println(sortedMap[keys[i]])
	// 	return sortedMap[keys[i]] < sortedMap[keys[j]]
	// })

	// for k, v := range sortedMap {
	// 	fmt.Printf("key: %s, value: %d\n", k, v)
	// }
	return hostAndTime
}

// mutexはコストが高く、十分なスケーラビリティを確保することができないため、Go1.9以降のバージョンでは、sync.Mapという並行安全なマップが追加されました[2]。
// func createTimeOnSiteMap(val TimeOnSite, hostAndTime *sync.Map, wg *sync.WaitGroup) {
// 	hostAndTime.Store(urlToHostName(val.Url), val.VisitDuration)
// 	defer wg.Done()
// 	// fmt.Println(hostAndTime)
// 	// 何故か↓のチャネルを追記したらhostAndTimeに値が入った
// 	// c <- true
// }

func urlToHostName(url string) string {
	rex := regexp.MustCompile("(http|https)://[\\w\\.-]+/")
	return rex.FindString(url)
}

// Mapの中で値の大きいバリューを持つキーを上位5つ探すメソッド
func searchTopFive(hostAndTime sync.Map) {
	topFiveValue := make([]int, 0, 5)
	topFiveKey := make([]string, 0, 5)
	hostAndTime.Range(func(key interface{}, value interface{}) bool {
		// キーからバリューを取り出し、値を比較してスライスに挿入
		// 先頭が一番大きい値が入る
		if len(topFiveValue) == 0 {
			topFiveKey = append(topFiveKey, key.(string))
			topFiveValue = append(topFiveValue, value.(int))
		}
		compareValue(&topFiveKey, &topFiveValue, key.(string), value.(int))
		return true
	})
	fmt.Println(topFiveValue)
	fmt.Println(topFiveKey)
}

// スライスの中で値の大きさが何番目かを比較
func compareValue(topFiveKey *[]string, topFiveValue *[]int, currentKey string, currentValue int) {
	for index, val := range *topFiveValue {
		if val < currentValue {
			// ポインタの値だけをコピー
			var copyTopFiveKey = *topFiveKey
			var copyTopFiveValue = *topFiveValue

			firstHalfKey := append(copyTopFiveKey[0:index], currentKey)
			secondHalfKey := copyTopFiveKey[index : len(copyTopFiveKey)-1]
			*topFiveKey = append(firstHalfKey, secondHalfKey...)

			firstHalfValue := append(copyTopFiveValue[0:index], currentValue)
			secondHalfValue := copyTopFiveValue[index : len(copyTopFiveValue)-1]
			*topFiveValue = append(firstHalfValue, secondHalfValue...)
		} 　
	}
}
