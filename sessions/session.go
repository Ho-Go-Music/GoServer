package sessions2

import (
	"github.com/gorilla/sessions"
	"sync"
)

var (
	Store *sessions.CookieStore
	once  sync.Once
)

func init() {
	once.Do(func() {
		Store = sessions.NewCookieStore([]byte("acaibird.com"))
	})
}
