package queries

import (
	store "github.com/jchenriquez/laundromat/store"
	models "github.com/jchenriquez/laundromat/types"
)

func GetUserByEmail(client *store.Client, email string) (models.User, error) {
	var user models.User
	statement := `Select * from "User" Where email=$1`
	row := client.Conn.QueryRow(statement, email)
	err := row.Scan(
		&user.FirstName, &user.LastName, &user.Email, &user.ID, &user.Country, &user.AdministrativeArea,
		&user.Thoroughfare, &user.Premise, &user.Locality, &user.Password, &user.PostalCode,
	)
	return user, err
}

func GetUserById(client *store.Client, id string) (models.User, error) {
	var user models.User
	statement := `Select * from "User" Where id=$1`
	row := client.Conn.QueryRow(statement, id)
	err := row.Scan(
		&user.FirstName, &user.LastName, &user.Email, &user.ID, &user.Country, &user.AdministrativeArea,
		&user.Thoroughfare, &user.Premise, &user.Locality, &user.Password, &user.PostalCode,
	)
	return user, err
}

func CreateUser(client *store.Client, user models.User) (models.User, error) {
	statement := `Insert 
						into "User" 
							(first_name, last_name, email, country, administrative_area, thoroughfare, premise, locality, postal_code, password)
							values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
						RETURNING id`

	row := client.Conn.QueryRow(
		statement, user.FirstName, user.LastName, user.Email, user.Country, user.AdministrativeArea, user.Thoroughfare,
		user.Premise, user.Locality, user.PostalCode, user.Password,
	)

	err := row.Scan(&user.ID)

	return user, err
}

func DeleteUser(client *store.Client, email string) error {
	stmt := `Delete From "User" where email=$1`
	_, err := client.Conn.Exec(stmt, email)
	return err
}

func UpdateUser(client *store.Client, user models.User) (models.User, error) {
	stmt := `UPDATE 
				"User" SET 
				first_name=$2, last_name=$3, email=$4, country=$5, administrative_area=$6, thoroughfare=$7, premise=$8, locality=$9, password=$10, postal_code=$11
			WHERE id=$1
			RETURNING *`

	row := client.Conn.QueryRow(stmt, user.ID, user.FirstName, user.LastName, user.Email, user.Country,
		user.AdministrativeArea, user.Thoroughfare, user.Premise, user.Locality, user.Password, user.PostalCode)
	err := row.Scan(&user.FirstName, &user.LastName, &user.Email, &user.ID, &user.Country, &user.AdministrativeArea,
		&user.Thoroughfare, &user.Premise, &user.Premise, &user.Locality, &user.Password, &user.PostalCode)

	return user, err
}
