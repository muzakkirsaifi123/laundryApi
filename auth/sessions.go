package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jchenriquez/laundromat/auth/cookies"
	jwt2 "github.com/jchenriquez/laundromat/auth/jwt"
	"github.com/jchenriquez/laundromat/store"
	"github.com/jchenriquez/laundromat/store/queries"
	"github.com/jchenriquez/laundromat/types"
)

type Session interface {
	Valid() bool
	AddFlash(key string, val interface{})
	GetFlash(key string) interface{}
}

type Manager interface {
	CreateSession(name string, userID int) Session
	SaveSession(session Session) error
	GetSession(name string) (Session, error)
	RefreshSession(session Session) error
	DeleteSession(name string) error
}

type SessionManager struct {
	jwtManager     jwt2.Manager
	cookiesManager cookies.Manager
	context        *gin.Context
	db             *store.Client
}

type SecureSession struct {
	UserID   int
	Name     string
	Flashes  map[string]interface{}
	InEffect bool
}

func NewSessionManager(c *gin.Context, db *store.Client) (Manager, error) {
	jwtManager, err := jwt2.StartJWTManager()
	if err != nil {
		return nil, err
	}

	cookiesManager := cookies.NewCookieManager()
	return SessionManager{
		jwtManager:     jwtManager,
		cookiesManager: cookiesManager,
		context:        c,
		db:             db,
	}, err
}

func (ss SecureSession) Valid() bool {
	return ss.InEffect
}

func (ss SecureSession) AddFlash(key string, val interface{}) {
	ss.Flashes[key] = val
}

func (ss SecureSession) GetFlash(key string) interface{} {
	return ss.Flashes[key]
}

func (sm SessionManager) CreateSession(name string, userID int) Session {
	return SecureSession{Name: name, InEffect: false, UserID: userID, Flashes: map[string]interface{}{"UserID": userID}}
}

func (sm SessionManager) SaveSession(session Session) error {
	securedSession := session.(SecureSession)
	accessToken, expirationTime, err := sm.jwtManager.GenerateAccessToken(securedSession.Flashes)
	refreshToken, _, err := sm.jwtManager.GenerateRefreshToken(securedSession.Flashes)

	if err != nil {
		return err
	}

	storeSession := types.Session{
		UserID:                securedSession.UserID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiration: expirationTime,
	}

	_, err = queries.CreateSession(sm.db, storeSession)
	if err != nil {
		return err
	}

	cookie, err := sm.cookiesManager.SessionCookie(securedSession.Flashes, securedSession.Name, expirationTime)

	if err != nil {
		return err
	}

	http.SetCookie(sm.context.Writer, &cookie)
	return nil
}

func (sm SessionManager) DeleteSession(name string) error {
	cookie, err := sm.context.Request.Cookie(name)
	if err != nil {
		return err
	}

	cookie.MaxAge = -1
	val, err := sm.cookiesManager.DecodeCookieValue(cookie)

	if err != nil {
		return err
	}

	userId, ok := val["UserID"]

	if !ok {
		return errors.New("malformed cookie")
	}

	err = queries.DeleteSessionForUser(sm.db, userId.(int))

	if err != nil {
		return err
	}
	http.SetCookie(sm.context.Writer, cookie)
	return nil
}

func (sm SessionManager) GetSession(name string) (Session, error) {
	cookie, err := sm.context.Request.Cookie(name)
	if err != nil {
		return nil, err
	}

	val, err := sm.cookiesManager.DecodeCookieValue(cookie)

	if err != nil {
		return nil, err
	}

	userId, ok := val["UserID"]
	securedSession := SecureSession{InEffect: false, Flashes: val, UserID: userId.(int)}

	if !ok {
		return securedSession, nil
	}

	typeSession, err := queries.GetSession(sm.db, userId.(int))

	securedSession.InEffect = sm.jwtManager.ValidateToken(
		typeSession.AccessToken, userId.(int),
	) || sm.jwtManager.ValidateToken(typeSession.RefreshToken, userId.(int))

	return securedSession, err
}

func (sm SessionManager) RefreshSession(session Session) error {
	securedSession := session.(SecureSession)
	accessToken, expirationTime, err := sm.jwtManager.GenerateAccessToken(securedSession.Flashes)
	refreshToken, _, err := sm.jwtManager.GenerateRefreshToken(securedSession.Flashes)

	if err != nil {
		return err
	}

	typedSession := types.Session{
		UserID:                securedSession.UserID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiration: expirationTime,
	}

	_, err = queries.CreateSession(sm.db, typedSession)

	return err
}
