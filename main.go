package main

import (
	"eleventh-learn/database"
	"eleventh-learn/router"
)

func main() {
	database.StartDB()
	r := router.StartApp()
	r.Run(":8000")
}
