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
type SiteInfo struct {
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
func GetData(config config.Config) []SiteInfo {
	// DB接続
	db, err := sql.Open("sqlite3", config.DbPath)
	if err != nil {
		log.Fatalf("Failed to connect db: %v", err)
	}
	defer db.Close()

	// クエリを実行
	// query := "SELECT title, url, last_visit_time FROM urls LIMIT 5"
	// 1ヶ月(30日)は2592000秒
	query := "SELECT urls.title, urls.url, visits.visit_duration FROM visits LEFT JOIN urls on visits.url = urls.id WHERE urls.last_visit_time >= (strftime('%s', 'now', '-2 years')+11644473600)*1000000 ORDER BY urls.last_visit_time ASC;"

	// query := "SELECT COUNT(title) FROM urls"
	rows, err := db.Query(query)
	var data []SiteInfo
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	// データ配列に格納
	for rows.Next() {
		var info SiteInfo
		err = rows.Scan(&info.Title, &info.Url, &info.VisitDuration)
		if err != nil {
			log.Fatalf("Failed to scan rows: %v", err)
		}
		data = append(data, info)
	}
	// for _, row := range data {
	// 	fmt.Println(row)
	// }
	return data

}

func GetLengthOfStay(datas []SiteInfo) ([]string, []int) {
	var hostAndTime2 = make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// データの個数分goroutineを実行するので、Addにはdatasの要素数を設定
	// wg.Done()が実行されるとAddが減っていく
	// DBからとってきた情報をkey: value = ホスト名: 滞在時間の合計に変換

	////////////////////////////////////
	// 1. まずキーに空文字？が入る場合があり、空文字が入るとRegex pattern unmatch!!が発生する
	// 2. 毎回結果が変わるのは、並列処理を行っている際、変数がキャプチャされる。処理の途中でバリューが変更されると正しい値が格納できない
	////////////////////////////////////
	var num int
	// 計測開始
	s := time.Now()
	for _, data := range datas {
		// バリューが上書きされてしまう
		wg.Add(1)
		go func(data SiteInfo) {
			// num += data.VisitDuration
			hostname := urlToHostName(data.Url)
			mu.Lock()

			if val, ok := hostAndTime2[hostname]; ok {
				hostAndTime2[hostname] = data.VisitDuration + val
			} else {
				hostAndTime2[hostname] = data.VisitDuration
			}

			// fmt.Println(urlToHostName(data.Url), data.VisitDuration)
			mu.Unlock()
			defer wg.Done()
		}(data)
	}
	// 経過時間を出力
	fmt.Printf("process time: %s\n", time.Since(s))
	fmt.Println(num)

	// Addが0になるまで待つ
	wg.Wait()

	// fmt.Println("Hello")
	// hostAndTime.Range(func(key interface{}, value interface{}) bool {
	// 	fmt.Printf("Key: %v(Type: %T) -> Value: %v(Type: %T)\n", key, key, value, value)
	// 	return true
	// })

	// topFiveKey, topFiveValue := getTopFive(hostAndTime)
	topFiveKey, topFiveValue := getTopFive(hostAndTime2)
	// fmt.Println(topFiveKey)
	// fmt.Println(topFiveValue)

	// キーのホスト名をサイト名に変換
	rex := regexp.MustCompile("([\\w-]+)\\.(com|co.jp|io)$")
	for index, hostname := range topFiveKey {
		match := rex.FindStringSubmatch(hostname)
		if match == nil {
			log.Fatalln("Regex pattern unmatch!!")
		}
		topFiveKey[index] = match[1]
	}

	// バリューのミリセカンドを時間に変換
	for index, milisecond := range topFiveValue {
		topFiveValue[index] = milisecond / (1000000 * 3600)
	}

	return topFiveKey, topFiveValue
}

// mutexはコストが高く、十分なスケーラビリティを確保することができないため、Go1.9以降のバージョンでは、sync.Mapという並行安全なマップが追加されました[2]。
// func createSiteInfoMap(val SiteInfo, hostAndTime *sync.Map, wg *sync.WaitGroup) {
// 	hostAndTime.Store(urlToHostName(val.Url), val.VisitDuration)
// 	defer wg.Done()
// 	// fmt.Println(hostAndTime)
// 	// 何故か↓のチャネルを追記したらhostAndTimeに値が入った
// 	// c <- true
// }

func urlToHostName(url string) string {
	rex := regexp.MustCompile("(http|https)://[\\w\\.-]+")
	return rex.FindString(url)
}

// Mapの中で値の大きいバリューを持つキーを上位5つ探すメソッド
func getTopFive(hostAndTime map[string]int) ([]string, []int) {
	topFiveKey := make([]string, 10, 10)
	topFiveValue := make([]int, 10, 10)
	for key, val := range hostAndTime {
		storeValueInOrder(&topFiveKey, &topFiveValue, key, val)
	}
	return topFiveKey, topFiveValue
}

// スライスの中で値の大きさが何番目かを比較
func storeValueInOrder(topFiveKey *[]string, topFiveValue *[]int, currentKey string, currentValue int) {
	for index, val := range *topFiveValue {
		if val < currentValue {
			copyTopFiveKey := make([]string, len(*topFiveKey))
			copy(copyTopFiveKey, *topFiveKey)

			// halfKey := copyTopFiveKey[:index]という記述にするとcopyTopFiveKeyのアドレスもコピーしてしまうのでcopyで値のみを代入
			halfKey := make([]string, index+1)
			copy(halfKey, copyTopFiveKey[:index])
			halfKey[index] = currentKey
			*topFiveKey = append(halfKey, copyTopFiveKey[index:len(copyTopFiveKey)-1]...)

			copyTopFiveValue := make([]int, len(*topFiveValue))
			copy(copyTopFiveValue, *topFiveValue)

			// halfValue := copyTopFiveValue[:index]という記述にするとcopyTopFiveValueのアドレスもコピーしてしまうのでcopyで値のみを代入
			halfValue := make([]int, index+1)
			copy(halfValue, copyTopFiveValue[:index])
			halfValue[index] = currentValue
			*topFiveValue = append(halfValue, copyTopFiveValue[index:len(copyTopFiveValue)-1]...)

			// fmt.Println(currentKey)
			// fmt.Println(len(*topFiveKey))
			// fmt.Println(topFiveKey, topFiveValue)
			break
		}
	}
}
