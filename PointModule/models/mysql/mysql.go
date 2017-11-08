package mysql

import (
	_ "github.com/Go-Sql-Driver/Mysql"
	"database/sql"
	"PointModule/logger"
)

var DB *sql.DB

func Init(dataSource string) {
	var err error
	DB, err = sql.Open("mysql", dataSource)
	if err != nil {
		logger.Fatal(err)
	}

	/*err = DB.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}*/
}