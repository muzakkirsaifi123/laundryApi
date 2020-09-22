package store

import (
	"fmt"
	"strings"
)

type Session struct {
	UID          int    `json:"uid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiration   string `json:"expiration"`
}

func (session *Session) Update(client *Client) error {
	return nil
}

func (session *Session) Get(client *Client) error {
	smt := fmt.Sprintf(`Select refresh_token, access_token from %s.public."Session" where uid = $1`, client.databaseName)
	row := client.db.QueryRow(client.Ctxt, smt, session.UID)
	return row.Scan(&session.RefreshToken, &session.AccessToken)
}

func (session *Session) Delete(client *Client) error {
	smt := fmt.Sprintf(`Delete from %s.public."Session" where uid = $1`, client.databaseName)
	_, err := client.db.Exec(client.Ctxt, smt, session.UID)
	return err
}

func (session *Session) Create(client *Client) error {
	currSession := Session{UID: session.UID}
	err := currSession.Get(client)

	if err != nil && !strings.Contains(err.Error(), "no rows") {
		return err
	}

	if len(currSession.AccessToken) > 0 || len(currSession.RefreshToken) > 0 {
		err = currSession.Delete(client)

		if err != nil {
			return err
		}
	}

	smt := fmt.Sprintf(`Insert into %s.public."Session" (uid, refresh_token, access_token) values ($1, $2, $3)`, client.databaseName)
	_, err = client.db.Exec(client.Ctxt, smt, session.UID, session.RefreshToken, session.AccessToken)
	return err
}
