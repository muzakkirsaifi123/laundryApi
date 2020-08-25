package store

import (
	"fmt"
	"strings"
)

type Order struct {
	ID         int `json:"id"`
	Lbs        int `json:"lbs"`
	UID        int `json:"user_id"`
	BusinessID int `json:"business_id"`
}

type Orders []Order

func (order *Order) Get(client *Client) error {
	statement := fmt.Sprintf(`Select lbs from %s.public."Order" where id = $1`, client.databaseName)
	row := client.db.QueryRow(client.Ctxt, statement, order.ID)

	return row.Scan(&order.Lbs)
}

func (order *Order) Create(client *Client) error {
	statement := fmt.Sprintf(`Insert into %s.public."Order" (lbs, user_id, business_id) values ($1, $2, $3)`, client.databaseName)
	_, err := client.db.Exec(client.Ctxt, statement, order.Lbs, order.UID, order.BusinessID)
	return err
}

func (order *Order) Delete(client *Client) error {
	statement := fmt.Sprintf(`Delete From %s.public."Order" where id=$1`, client.databaseName)
	_, err := client.db.Exec(client.Ctxt, statement, order.ID)
	return err
}

func (order *Order) Update(client *Client) error {
	statement := fmt.Sprintf(`UPDATE %s.public."Order" SET lbs=$1 where id=$2`, client.databaseName)
	_, err := client.db.Exec(client.Ctxt, statement, order.Lbs, order.ID)
	return err
}

func (orders *Orders) Get(client *Client, query Query) error {
	statement := fmt.Sprintf(`Select * from %s.public."Order"`, client.databaseName)
	builder := strings.Builder{}
	count := 0
	keys := []string{"id", "lbs", "user_id", "business_id"}
	for _, key := range keys {
		val := query.Get(key)

		if len(val) > 0 {
			if count == 0 {
				builder.WriteString(fmt.Sprintf(`%s='%v'`, key, query.Get(key)))
			} else {
				builder.WriteString(fmt.Sprintf(` AND %s='%v'`, key, query.Get(key)))
			}
			count++
		}
	}

	queryStatement := fmt.Sprintf("%s", statement)

	if builder.Len() > 0 {
		queryStatement = fmt.Sprintf("%s where %s", queryStatement, builder.String())
	}
	rows, err := client.db.Query(client.Ctxt, queryStatement)

	if err != nil {
		return err
	}

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.Lbs, &order.UID, &order.BusinessID)

		if err != nil {
			return err
		}

		*orders = append(*orders, order)
	}

	return nil
}
