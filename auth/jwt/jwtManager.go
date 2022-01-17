package jwt

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	config "github.com/spf13/viper"
)

const REFRESH_TOKEN_EXP_HOURS_EXTENSION = 5

type Manager struct {
	Secret            string
	ExpirationHours   int64
	ExpirationMinutes int64
}

type ManagerClaim struct {
	Email string
	jwt.StandardClaims
}

var jwtManager Manager

func StartJWTManager() (Manager, error) {
	jwtInfo := config.GetStringMapString("jwt")
	secret := jwtInfo["secret"]
	expirationHours, err := strconv.Atoi(jwtInfo["expirationhours"])
	if err != nil {
		return Manager{}, err
	}
	expirationMinutes, err := strconv.Atoi(jwtInfo["expirationminutes"])
	if err != nil {
		return Manager{}, err
	}

	jwtManager = Manager{secret, int64(expirationHours), int64(expirationMinutes)}
	return jwtManager, nil
}

func (mJwt Manager) ValidateToken(tokenStr string, userId int) bool {
	jwtTokenManager := jwtauth.New(jwt.SigningMethodHS256.Alg(), []byte(mJwt.Secret), nil)
	token, err := jwtTokenManager.Decode(tokenStr)

	if err != nil {
		return false
	}

	if token.Expiration().UnixNano() < time.Now().UnixNano() {
		return false
	}

	id, ok := token.PrivateClaims()["UserID"]

	if !ok {
		return false
	}

	switch v := id.(type) {
	case int:
		return v == userId
	case float64:
		return int(v) == userId
	case string:
		claimId, err := strconv.Atoi(v)
		if err != nil {
			return false
		}

		return claimId == userId
	default:
		return false
	}
}

func (mJwt Manager) GenerateAccessToken(claims jwt.MapClaims) (string, time.Time, error) {
	jwtTokenManager := jwtauth.New(jwt.SigningMethodHS256.Alg(), []byte(mJwt.Secret), nil)

	t := time.Now().Add(time.Hour*time.Duration(mJwt.ExpirationHours) + time.Minute*time.Duration(mJwt.ExpirationMinutes))
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiry(claims, t)
	_, tokenStr, err := jwtTokenManager.Encode(claims)
	return tokenStr, t, err
}

func (mJwt Manager) GenerateRefreshToken(claims jwt.MapClaims) (string, string, error) {
	jwtTokenManager := jwtauth.New(jwt.SigningMethodHS256.Alg(), []byte(mJwt.Secret), nil)

	t := time.Now().Add(time.Hour*time.Duration(mJwt.ExpirationHours+REFRESH_TOKEN_EXP_HOURS_EXTENSION) + time.Minute*time.Duration(mJwt.ExpirationMinutes))
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiry(claims, t)
	_, tokenStr, err := jwtTokenManager.Encode(claims)
	return tokenStr, t.String(), err
}
