package queries

import (
	"github.com/jchenriquez/laundromat/store"
	"github.com/jchenriquez/laundromat/types"
)

func GetSession(client *store.Client, userId int) (types.Session, error) {
	var ret types.Session
	smt := `Select refresh_token, access_token from "Session" where user_id = $1`
	row := client.Conn.QueryRow(smt, userId)
	err := row.Scan(&ret.RefreshToken, &ret.AccessToken)
	return ret, err
}

func DeleteSessionForUser(client *store.Client, id int) error {
	smt := `Delete from "Session" where user_id = $1`
	_, err := client.Conn.Exec(smt, id)
	return err
}

func DeleteSessions(client *store.Client) error {
	smt := `Delete from "Session" where true`
	_, err := client.Conn.Exec(smt)
	return err
}

func CreateSession(client *store.Client, session types.Session) (types.Session, error) {
	smt := `Insert into "Session" (user_id, refresh_token, access_token, access_token_expiration) values ($1, $2, $3, $4)
			On CONFLICT(user_id)
			Do 
				Update Set refresh_token=Excluded.refresh_token, access_token=Excluded.access_token, access_token_expiration=Excluded.access_token_expiration
			RETURNING id`
	conn := client.Conn.QueryRow(
		smt, session.UserID, session.RefreshToken, session.AccessToken, session.AccessTokenExpiration,
	)
	err := conn.Scan(&session.ID)
	return session, err
}
