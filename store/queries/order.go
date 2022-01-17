package queries

import (
	"github.com/jchenriquez/laundromat/store"
	"github.com/jchenriquez/laundromat/types"
)

func GetOrder(client *store.Client, orderId string) (types.Order, error) {
	var order types.Order
	statement := `Select * from "Order" where id = $1`
	row := client.Conn.QueryRow(statement, orderId)

	err := row.Scan(
		&order.ID, &order.UserID, &order.PreferredDetergent, &order.UseBleach, &order.PreferredSoftener, &order.WeightLBS,
		&order.PickupCharge, &order.DropCharge, &order.ServiceCharge,
	)
	return order, err
}

func CreateOrder(client *store.Client, order types.Order) (types.Order, error) {
	stmt := `Insert 
				into 
					"Order" (user_id, business_id, preferred_detergent, use_bleach, preferred_softener, weight_lbs, dropoff_charge, service_charge) 
					values ($1, $2, $3, $4, $5, $6, $7, $8)
					RETURNING id`

	row := client.Conn.QueryRow(stmt, order.UserID, order.BusinessID, order.PreferredDetergent, order.UseBleach, order.PreferredSoftener,
		order.WeightLBS, order.DropCharge, order.ServiceCharge,
	)
	err := row.Scan(&order.ID)
	return order, err
}

func DeleteOrder(client *store.Client, order types.Order) error {
	stmt := `Delete From Order where id=$1`
	_, err := client.Conn.Exec(stmt, order.ID)
	return err
}

func UpdateOrder(client *store.Client, order types.Order) (types.Order, error) {
	stmt := `UPDATE "Order" 
		SET preferred_detergent=$2, user_bleach=$3, preferred_softener=$4, weight_lbs=$5, pickup_charge=$6, dropoff_charge=$7, service_charge=$8
		WHERE id=$1
        RETURNING preferred_detergent, user_bleach, preferred_softener, weight_lbs, pickup_charge, dropoff_charge, service_charge`

	conn := client.Conn.QueryRow(stmt, order.ID, order.PreferredDetergent, order.UseBleach, order.PreferredSoftener, order.WeightLBS, order.PickupCharge,
		order.DropCharge, order.ServiceCharge)
	err := conn.Scan(&order.PreferredDetergent, &order.UseBleach, &order.PreferredSoftener, &order.WeightLBS, &order.DropCharge, &order.ServiceCharge)
	return order, err
}

func GetOrders(client *store.Client, query map[string]interface{}) ([]types.Order, error) {

	stmt := `Select * from "Order" where id=$1 OR user_id=$2 OR business_id=$3`
	rows, err := client.Conn.Query(stmt, query["id"], query["user_id"], query["business_id"])

	if err != nil {
		return []types.Order{}, err
	}
	defer rows.Close()
	orders := make([]types.Order, 0)

	for rows.Next() {
		var order types.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.PreferredDetergent, &order.UseBleach, &order.PreferredSoftener, &order.WeightLBS,
			&order.PickupCharge, &order.DropCharge, &order.ServiceCharge, &order.BusinessID,
		)

		if err != nil {
			return orders, err
		}

		orders = append(orders, order)
	}

	return orders, err
}
