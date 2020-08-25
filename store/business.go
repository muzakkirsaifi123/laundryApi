package store

import (
	"fmt"
	"strings"
)

type Business struct {
	ID           int    `json:"id"`
	BusinessName string `json:"business_name"`
	Address      string `json:"address"`
	Rating       int    `json:"rating"`
}

type Businesses []Business

func (business *Business) Get(client *Client) error {
	statement := fmt.Sprintf(`Select business_name, address, rating from %s.public."Business" where id = $1`, client.databaseName)
	row := client.db.QueryRow(client.Ctxt, statement, business.ID)

	return row.Scan(&business.BusinessName, &business.Address, &business.Rating)
}

func (business *Business) Create(client *Client) error {
	statement := fmt.Sprintf(`Insert into %s.public."Business" (business_name, address, rating) values ($1, $2, $3)`, client.databaseName)
	_, err := client.db.Exec(client.Ctxt, statement, business.BusinessName, business.Address, business.Rating)
	return err
}

func (business *Business) Delete(client *Client) error {
	statement := fmt.Sprintf(`Delete From %s.public."Business" where id=$1`, client.databaseName)
	_, err := client.db.Exec(client.Ctxt, statement, business.ID)
	return err
}

func (business *Business) Update(client *Client) error {
	statement := fmt.Sprintf(`UPDATE %s.public."Business" SET business_name=$1, SET address=$2 where rating=$3`, client.databaseName)
	_, err := client.db.Exec(client.Ctxt, statement, business.BusinessName, business.Address, business.Rating)
	return err
}

func (businesses *Businesses) Get(client *Client, query Query) error {
	statement := fmt.Sprintf(`Select * from %s.public."Business"`, client.databaseName)
	builder := strings.Builder{}
	keys := []string{"business_name", "address", "rating"}
	count := 0
	for _, key := range keys {
		val := query.Get(key)

		if len(val) > 0 {
			if count == 0 {
				builder.WriteString(fmt.Sprintf(`%s='%v'`, key, val))
			} else {
				builder.WriteString(fmt.Sprintf(` AND %s='%v'`, key, val))
			}
			count++
		}
	}

	queryStatement := fmt.Sprintf("%s", statement)

	if builder.Len() > 0 {
		queryStatement = fmt.Sprintf("%s where %s", queryStatement, builder.String())
	}

	fmt.Printf("queryStatement %s\n", queryStatement)
	rows, err := client.db.Query(client.Ctxt, queryStatement)

	if err != nil {
		return err
	}

	for rows.Next() {
		var business Business

		err := rows.Scan(&business.ID, &business.BusinessName, &business.Address, &business.Rating)

		if err != nil {
			return err
		}

		*businesses = append(*businesses, business)
	}

	return nil
}
