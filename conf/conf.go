package conf

import (
	"fmt"
)

var (
	mysqlHost = "127.0.0.1:3306"
	mysqlUsername = "u_baopo"
	mysqlPassword = "123456"
	mysqlDb = "db_point"

	DataSource = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", mysqlUsername, mysqlPassword, mysqlHost, mysqlDb)

	LogPath = "/Users/didi/baopo/go/logs/PointModule.log"
)
