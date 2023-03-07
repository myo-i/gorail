package main

import (
	"gorail/db"
	"gorail/util"
)

func main() {
	db.GetData()
	util.Load()
	// user.RunApp()
}
