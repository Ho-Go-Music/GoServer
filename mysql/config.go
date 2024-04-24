package mysql

import (
	diylog "acaibird.com/log"
	"acaibird.com/tools"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DatabaseUrl = tools.Conf.MySQLServer.User + ":" + tools.Conf.MySQLServer.Password +
		"@" + "tcp" +
		"(" + tools.Conf.MySQLServer.Host + ":" + tools.Conf.MySQLServer.Port + ")" +
		"/" + tools.Conf.MySQLServer.Database +
		"?" + "charset=" + tools.Conf.MySQLServer.Charset +
		"&" + "parseTime=" + tools.Conf.MySQLServer.ParseTime +
		"&" + "loc=" + tools.Conf.MySQLServer.Loc
)

func Newdb() (*sql.DB, error) {
	db, err := sql.Open("mysql", DatabaseUrl)
	if err != nil {
		diylog.Sugar.Errorln(err)
	}
	return db, err
}
