package httpfunc

import (
	diylog "acaibird.com/log"
	"acaibird.com/mysql"
	"acaibird.com/mysql/table"
	"acaibird.com/request_body"
	"database/sql"
	"encoding/json"
	"net/http"
)

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	logger := diylog.NewLogger()
	defer func() {
		err := logger.Sync()
		if err != nil {
			return
		}
	}()
	// request method check
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user request_body.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := mysql.Newdb()
	if err != nil {
		logger.Errorln("Database error", err)
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return

	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Errorln(err)
		}
	}(db)

	if userExit := table.IsUserExists(user.Username, db); !userExit {
		http.Error(w, "NotFound", http.StatusNotFound)
		return
	}
	if passwordCorrect := table.VerifyPassword(user.Username, user.Password, db); !passwordCorrect {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if isActive := table.IsActive(user.Username, db); !isActive {
		http.Error(w, "NotActive", http.StatusForbidden)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
