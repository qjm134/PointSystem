package conf

import (
	"fmt"
)

var (
	mysqlHost = "127.0.0.1:3306"
	mysqlUsername = "xxxx"
	mysqlPassword = "xxxx"
	mysqlDb = "xxxx"

	DataSource = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", mysqlUsername, mysqlPassword, mysqlHost, mysqlDb)

	LogPath = "xxxx/logs/PointModule.log"
)
