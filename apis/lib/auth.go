package lib

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/jchenriquez/laundromat/store"
	config "github.com/spf13/viper"
	"net/http"
	"time"
)

func Authorizer(clientDB *store.Client, sessionToken func(session store.Session) (*jwt.Token, error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token, claims, err := jwtauth.FromContext(request.Context())

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			session := store.Session{UID: int(claims["UID"].(float64))}
			err = session.Get(clientDB)

			if err != nil {
				http.Error(writer, "not authorized", http.StatusUnauthorized)
				return
			}

			dbToken, err := sessionToken(session)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			if dbToken.Raw != token.Raw {
				http.Error(writer, "not authorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func AccessAuthorizer(clientDB *store.Client) func(next http.Handler) http.Handler {
	jwtSecret := config.GetString("jwt_secret")
	jwtTokenManager := jwtauth.New(jwt.SigningMethodHS256.Alg(), []byte(jwtSecret), nil)

	return Authorizer(clientDB, func(session store.Session) (*jwt.Token, error) {
		return jwtTokenManager.Decode(session.AccessToken)
	})
}

func RefresherAuthorizer(clientDB *store.Client) func(next http.Handler) http.Handler {
	jwtSecret := config.GetString("jwt_refresher")
	jwtTokenManager := jwtauth.New(jwt.SigningMethodHS256.Alg(), []byte(jwtSecret), nil)

	return Authorizer(clientDB, func(session store.Session) (*jwt.Token, error) {
		return jwtTokenManager.Decode(session.RefreshToken)
	})
}

func Auth(client *store.Client) func(handler http.Handler) http.Handler {
	jwtSecret := config.GetString("jwt_secret")
	jwtTokenManager := jwtauth.New(jwt.SigningMethodHS256.Alg(), []byte(jwtSecret), nil)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Method == http.MethodPost && request.URL.Path == "/user" {
				next.ServeHTTP(writer, request)
			} else {
				chi.Chain(jwtauth.Verifier(jwtTokenManager), jwtauth.Authenticator, AccessAuthorizer(client)).Handler(next).
					ServeHTTP(writer, request)
			}
		})
	}
}

func RefreshAuth(client *store.Client) func(next http.Handler) http.Handler {
	jwtSecret := config.GetString("jwt_refresher")
	jwtTokenManager := jwtauth.New(jwt.SigningMethodHS256.Alg(), []byte(jwtSecret), nil)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			chi.Chain(jwtauth.Verifier(jwtTokenManager), jwtauth.Authenticator, RefresherAuthorizer(client)).Handler(next).
				ServeHTTP(writer, request)
		})
	}
}

func GenerateAccessToken(user store.User) (string, string, error) {
	jwtSecret := config.GetString("jwt_secret")
	jwtTokenManager := jwtauth.New(jwt.SigningMethodHS256.Alg(), []byte(jwtSecret), nil)

	claim := jwt.MapClaims{
		"UID": user.UID,
	}
	t := time.Now().Add(time.Minute * 5)
	jwtauth.SetIssuedNow(claim)
	jwtauth.SetExpiry(claim, t)
	_, tokenStr, err := jwtTokenManager.Encode(claim)
	return tokenStr, t.String(), err
}

func GenerateRefreshToken(user store.User) (string, string, error) {
	jwtSecret := config.GetString("jwt_refresher")
	jwtTokenManager := jwtauth.New(jwt.SigningMethodHS256.Alg(), []byte(jwtSecret), nil)

	claim := jwt.MapClaims{
		"UID": user.UID,
	}
	t := time.Now().Add(time.Minute * 15)
	jwtauth.SetIssuedNow(claim)
	jwtauth.SetExpiry(claim, t)
	_, tokenStr, err := jwtTokenManager.Encode(claim)
	return tokenStr, t.String(), err
}
