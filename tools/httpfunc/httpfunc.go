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
	"fmt"
	"github.com/gorilla/sessions"
	redis2 "github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	// log
	logger := diylog.Sugar
	defer func() {
		err := logger.Sync()
		if err != nil {
			return
		}
	}()
	// redis
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	redisClient := redis.GetRedisClient()
	defer func(redisClient *redis2.Client) {
		err := redisClient.Close()
		if err != nil {

		}
	}(redisClient)
	// sql
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
	// request method check
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// parse request body
	var user request_body.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		diylog.Sugar.Errorln("json decode throw a error ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("repeat login")
		return
	}
	//// repeat login check
	//val, _ := redisClient.Get(ctx, user.Username).Result()
	//if val != "" {
	//	http.Error(w, "Repeat login", http.StatusConflict)
	//	return
	//}
	// user dose not exit
	if userExit := table.IsUserExists(user.Username, db); !userExit {
		http.Error(w, "NotFound", http.StatusNotFound)
		return
	}
	// password error
	if passwordCorrect := table.VerifyPassword(user.Username, user.Password, db); !passwordCorrect {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// user status is disable
	if isActive := table.IsActive(user.Username, db); !isActive {
		http.Error(w, "NotActive", http.StatusForbidden)
		return
	}
	// 编码器（发送器）和解码器（接收器）之间进行二进制数据流的发送，
	// 一般用来传递远端程序调用的参数和结果，比如net/rpc包就有用到这个
	gob.Register(time.Time{})
	// session
	session, err := sessions2.Store.Get(r, "session")
	// save session data to cookie
	session.Values[user.Username] = time.Now()
	session.Options = &sessions.Options{
		MaxAge:   60 * 60 * 2,
		Domain:   "localhost",
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: false,
	}
	err = session.Save(r, w)
	if err != nil {
		diylog.Sugar.Errorln("session save throw a error ", err)
	}
	// save user logged time in redis
	_, err = redisClient.Set(ctx, user.Username, time.Now(), 60*time.Minute).Result()
	if err != nil {
		diylog.Sugar.Errorln("redis save online user info throw a error ", err)
	}
	diylog.Sugar.Infoln("redis save online user log in time successfully\n")
	// response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
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
	// get user's session
	session, err := sessions2.Store.Get(r, "session")
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}
	// delete user's session,MaxAge: -1 means delete the session
	Domain := strings.Join(strings.Split(r.Host, ":"), "")
	session.Options = &sessions.Options{
		Path:     "/",
		Domain:   Domain,
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	delete(session.Values, username)
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}
	// delete user's online info in redis
	redisClient := redis.GetRedisClient()
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	_, err = redisClient.Del(ctx, username).Result()
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}
	// response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("success"))
	if err != nil {
		return
	}
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessions2.Store.Get(r, "session")
	if err != nil {
		println(err.Error())
	}
	if v, ok := session.Values["root"]; !ok {
		print("no root")
	} else {
		diylog.Sugar.Infoln(v)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("success"))
	if err != nil {
		return
	}
}
