package cookies

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
)

type Manager interface {
	SessionCookie(values map[string]interface{}, name string, expires time.Time, maxAgeSeconds int) (http.Cookie, error)
	DecodeCookieValue(cookie *http.Cookie) (map[string]interface{}, error)
}

type cookieGenerator struct {
	hashKey  []byte
	blockKey []byte
}

var once sync.Once
var cookieStore cookieGenerator

func NewCookieManager() Manager {
	once.Do(
		func() {
			cookieStore = cookieGenerator{securecookie.GenerateRandomKey(16), securecookie.GenerateRandomKey(16)}
		},
	)

	return cookieStore
}

func (cg cookieGenerator) SessionCookie(
	values map[string]interface{}, name string, expires time.Time, maxAgeSeconds int,
) (
	http.Cookie, error,
) {
	sc := securecookie.New(cg.hashKey, cg.blockKey)
	encoded, err := sc.Encode(
		name, values,
	)

	if err != nil {
		return http.Cookie{}, err
	}

	return http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  expires,
		MaxAge:   maxAgeSeconds,
	}, nil
}

func (cg cookieGenerator) DecodeCookieValue(cookie *http.Cookie) (map[string]interface{}, error) {
	sc := securecookie.New(cg.hashKey, cg.blockKey)
	value := make(map[string]interface{})
	err := sc.Decode(cookie.Name, cookie.Value, &value)
	return value, err
}
