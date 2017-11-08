package main

import (
	"net/http"
	"PointModule/conf"
	"PointModule/models/mysql"
	"PointModule/controllers"
	"PointModule/logger"
)

func main() {
	logger.Init()
	mysql.Init(conf.DataSource)

	webService()
}

func webService() {
	http.HandleFunc("/getTotalPoints", controllers.GetTotalPoints)
	http.HandleFunc("/addPoints", controllers.AddPoints)
	http.HandleFunc("/deductPoints", controllers.DeductPoints)
	http.HandleFunc("/getPointsRecords", controllers.GetPointsRecords)

	err := http.ListenAndServe(":8017", nil)
	if err != nil {
		logger.Fatal(err)
	}
}
