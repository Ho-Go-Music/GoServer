package table

import (
	diylog "acaibird.com/log"
	"database/sql"
)

type Users struct {
	ID       int
	name     string
	password string
	email    string
	active   int
}

func IsUserExists(username string, db *sql.DB) bool {
	logger := diylog.NewLogger()
	defer func() {
		err := logger.Sync()
		if err != nil {
			return
		}
	}()

	smtp, err := db.Prepare("SELECT * FROM users WHERE name = ?")
	if err != nil {
		logger.Errorln("prepare sql statement exception: ", err)
		return false
	}

	row, err := smtp.Query(username)
	defer func() {
		err := row.Close()
		if err != nil {
			return
		}
	}()

	for row.Next() {
		return true
	}
	if err != nil {
		logger.Errorln("execute sql statement exception: ", err)
		return false
	}
	return false
}

func IsActive(username string, db *sql.DB) bool {
	logger := diylog.NewLogger()
	defer func() {
		err := logger.Sync()
		if err != nil {
			return
		}
	}()

	smtp, err := db.Prepare("SELECT active FROM users WHERE name = ?")
	if err != nil {
		logger.Errorln("prepare sql statement exception: ", err)
		return false
	}

	row, err := smtp.Query(username)
	defer func() {
		err := row.Close()
		if err != nil {
			return
		}
	}()

	var active int
	for row.Next() {
		err := row.Scan(&active)
		if err != nil {
			logger.Errorln("scan sql statement exception: ", err)
			return false
		}
	}

	if active == 1 {
		return true
	}
	return false
}

func VerifyPassword(username, password string, db *sql.DB) bool {
	logger := diylog.NewLogger()
	defer func() {
		err := logger.Sync()
		if err != nil {
			return
		}
	}()

	smtp, err := db.Prepare("SELECT password FROM users WHERE name = ?")
	if err != nil {
		logger.Errorln("prepare sql statement exception: ", err)
		return false
	}

	row, err := smtp.Query(username)
	defer func() {
		err := row.Close()
		if err != nil {
			return
		}
	}()

	var realPassword string
	for row.Next() {
		err := row.Scan(&realPassword)
		if err != nil {
			logger.Errorln("scan sql statement exception: ", err)
			return false
		}
	}

	if realPassword == password {
		return true
	}
	return false
}
