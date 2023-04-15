package main

import (
	"fmt"
	"gorail/db"
)

func main() {
	data := db.GetData()
	result := db.CalcTimeOnSite(data)
	fmt.Println(result.Load("s"))

	// config.Load()
	// user.RunApp()
}
