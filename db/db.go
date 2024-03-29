package db

import (
	"database/sql"
	"gorail/config"
	"log"
	"regexp"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// DBからデータを取得した際の構造体
type SiteInfo struct {
	Title         string
	Url           string
	VisitDuration int
}

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
	query := "SELECT urls.title, urls.url, visits.visit_duration FROM visits LEFT JOIN urls on visits.url = urls.id WHERE urls.last_visit_time >= (strftime('%s', 'now', '-2 months')+11644473600)*1000000 ORDER BY urls.last_visit_time ASC;"

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
	var hostAndTime = make(map[string]int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	////////////////////////////////////
	// 1. まずキーに空文字？が入る場合があり、空文字が入るとRegex pattern unmatch!!が発生する
	// 2. 毎回結果が変わるのは、並列処理を行っている際、変数がキャプチャされる。処理の途中でバリューが変更されると正しい値が格納できない
	////////////////////////////////////
	for _, data := range datas {
		// バリューが上書きされてしまう
		wg.Add(1)
		go func(data SiteInfo) {
			// num += data.VisitDuration
			hostname := urlToHostName(data.Url)
			mu.Lock()

			if val, ok := hostAndTime[hostname]; ok {
				hostAndTime[hostname] = data.VisitDuration + val
			} else {
				hostAndTime[hostname] = data.VisitDuration
			}

			mu.Unlock()
			defer wg.Done()
		}(data)
	}

	// Addが0になるまで待つ
	wg.Wait()

	topFiveKey, topFiveValue := getTopFive(hostAndTime)

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

			halfKey := make([]string, index+1)
			copy(halfKey, copyTopFiveKey[:index])
			halfKey[index] = currentKey
			*topFiveKey = append(halfKey, copyTopFiveKey[index:len(copyTopFiveKey)-1]...)

			copyTopFiveValue := make([]int, len(*topFiveValue))
			copy(copyTopFiveValue, *topFiveValue)

			halfValue := make([]int, index+1)
			copy(halfValue, copyTopFiveValue[:index])
			halfValue[index] = currentValue
			*topFiveValue = append(halfValue, copyTopFiveValue[index:len(copyTopFiveValue)-1]...)

			break
		}
	}
}
