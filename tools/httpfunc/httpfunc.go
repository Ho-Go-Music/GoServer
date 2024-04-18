package httpfunc

import (
	diylog "acaibird.com/log"
	"acaibird.com/mysql"
	"acaibird.com/mysql/table"
	"acaibird.com/redis"
	"acaibird.com/request_body"
	sessions2 "acaibird.com/sessions"
	"context"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	logger := diylog.Sugar
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

	// 编码器（发送器）和解码器（接收器）之间进行二进制数据流的发送，
	// 一般用来传递远端程序调用的参数和结果，比如net/rpc包就有用到这个
	gob.Register(time.Time{})
	session, err := sessions2.Store.Get(r, "session")
	session.Values[user.Username] = time.Now()
	session.Options = &sessions.Options{
		MaxAge: 60 * 60 * 2,
	}
	err = session.Save(r, w)
	if err != nil {
		diylog.Sugar.Errorln("session save throw a error ", err)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	redisClient := redis.GetRedisClient()
	_, err = redisClient.Set(ctx, user.Username, time.Now(), 60*time.Minute).Result()
	if err != nil {
		diylog.Sugar.Errorln("redis save online user info throw a error ", err)
	}
	diylog.Sugar.Infoln("redis save online user info successfully")

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("success"))
	if err != nil {
		return
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	session, err := sessions2.Store.Get(r, name)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}
	if session.IsNew {
		session.Save(r, w)
		println("session is new")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, " + name))
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	session, err := sessions2.Store.Get(r, "session")
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}
	session.Options = &sessions.Options{
		MaxAge: -1,
	}
	delete(session.Values, username)
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	redisClient := redis.GetRedisClient()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_, err = redisClient.Del(ctx, username).Result()
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("success"))
	if err != nil {
		return
	}
}
