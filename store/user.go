package store

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID      int    `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (user *User) Get(client *Client) error {
	statement := fmt.Sprintf(`Select * from %s.public."User" where uid = $1 or email = $2`, client.databaseName)
	row := client.db.QueryRow(client.Ctxt, statement, user.UID, user.Email)
	return row.Scan(&user.Email, &user.UID, &user.Password)
}

func (user *User) Create(client *Client) error {
	statement := fmt.Sprintf(`Insert into %s.public."User" (email, password) values ($1, $2)`, client.databaseName)
	password := user.Password

	if len(password) == 0 || len(user.Email) == 0 {
		return errors.New("user must have email and password")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = client.db.Exec(client.Ctxt, statement, user.Email, string(hashed))
	return err
}

func (user *User) Delete(client *Client) error {
	statement := fmt.Sprintf(`Delete From %s.public."User" where uid=$1`, client.databaseName)
	_, err := client.db.Exec(client.Ctxt, statement, user.UID)
	return err
}

func (user *User) Update(client *Client) error {
	statement := fmt.Sprintf(`UPDATE %s.public."User" SET email=$1, SET address=$2 SET password=$3 where uid=$3`, client.databaseName)
	_, err := client.db.Exec(client.Ctxt, statement, user.Email, user.UID, user.Password)
	return err
}
