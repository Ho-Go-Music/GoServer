package mysql

import (
	diylog "acaibird.com/log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DatabaseUrl = "root:775028@tcp(127.0.0.1:3306)/WebMusic?charset=utf8mb4&parseTime=True&loc=Local"
)

func Newdb() (*sql.DB, error) {
	db, err := sql.Open("mysql", DatabaseUrl)
	if err != nil {
		diylog.NewLogger().Errorln(err)
	}
	return db, err
}
