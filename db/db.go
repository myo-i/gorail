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

// とりあえず3か月分くらいのデータを（できれば滞在時間が長い順で）取得する
// 取得したデータを1日、1週間、1か月とかでさらに範囲を絞って使う

type History interface {
	GetHistory()
}

// DBからデータを取得
func GetData() []SiteInfo {
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
	query := "SELECT urls.title, urls.url, visits.visit_duration FROM visits LEFT JOIN urls on visits.url = urls.id WHERE urls.last_visit_time >= (strftime('%s', 'now', '-1 months')+11644473600)*1000000 ORDER BY urls.last_visit_time ASC;"

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
	var hostAndTime = sync.Map{}
	var wg sync.WaitGroup

	// データの個数分goroutineを実行するので、Addにはdatasの要素数を設定
	// wg.Done()が実行されるとAddが減っていく
	wg.Add(len(datas))
	// DBからとってきた情報をkey: value = ホスト名: 滞在時間の合計に変換
	for _, data := range datas {
		// バリューが上書きされてしまう
		go func(data SiteInfo) {
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

	topFiveKey, topFiveValue := getTopFive(hostAndTime)

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
	rex := regexp.MustCompile("(http|https)://[\\w\\.-]+/")
	return rex.FindString(url)
}

// Mapの中で値の大きいバリューを持つキーを上位5つ探すメソッド
func getTopFive(hostAndTime sync.Map) ([]string, []int) {
	topFiveKey := make([]string, 5, 5)
	topFiveValue := make([]int, 5, 5)
	hostAndTime.Range(func(key interface{}, value interface{}) bool {
		// メソッドにする意味あんまないかも
		storeValueInOrder(&topFiveKey, &topFiveValue, key.(string), value.(int))
		// fmt.Println(topFiveKey)
		// fmt.Println(topFiveValue)
		// fmt.Println("----------------------------------------------------")
		return true
	})
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
			break
		}
	}
}
