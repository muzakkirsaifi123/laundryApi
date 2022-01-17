package queries

import (
	"strconv"

	"github.com/jchenriquez/laundromat/store"
	"github.com/jchenriquez/laundromat/types"
)

func GetBusiness(client *store.Client, id int) (types.Business, error) {
	var business types.Business
	stmt := `Select * from "Business" where id = $1`
	conn := client.Conn.QueryRow(stmt, id)
	err := conn.Scan(
		&business.OrganizationName, &business.Country, &business.AdministrativeArea, &business.Locality,
		&business.PostalCode,
		&business.Premise, &business.Thoroughfare, &business.ID, &business.Rating, &business.PricePerPound,
	)
	return business, err
}

func CreateBusiness(client *store.Client, business types.Business) (types.Business, error) {
	stmt := `Insert 
				into "Business" (organization_name, country, administrative_area, locality, postal_code, premise, thoroughfare, id, rating, price_per_pound) 
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				RETURNING *`
	conn := client.Conn.QueryRow(
		stmt, business.OrganizationName, business.Country, business.AdministrativeArea, business.Locality,
		business.PostalCode, business.Premise, business.Thoroughfare, business.ID, business.Rating,
		business.PricePerPound,
	)
	err := conn.Scan(
		&business.OrganizationName, &business.Country, &business.AdministrativeArea, &business.Locality,
		&business.PostalCode, &business.Premise, &business.Thoroughfare, &business.ID, &business.Rating,
		&business.PricePerPound,
	)
	return business, err
}

func DeleteBusiness(client *store.Client, id int) error {
	stmt := `Delete From "Business" where id=$1`
	_, err := client.Conn.Exec(stmt, id)
	return err
}

func UpdateBusiness(client *store.Client, business types.Business) (types.Business, error) {
	stmt := `UPDATE "Business" 
				SET organization_name, country, administrative_area, locality, postal_code, premise, thoroughfare, rating, price_per_pound
				WHERE id=$4
				RETURNING *`

	conn := client.Conn.QueryRow(stmt, business.ID)
	err := conn.Scan(
		&business.OrganizationName, &business.Country, &business.AdministrativeArea, &business.Locality,
		&business.PostalCode, &business.Premise, &business.Thoroughfare,
		&business.Rating, &business.PricePerPound,
	)
	return business, err
}

func GetBusinesses(client *store.Client, query map[string]interface{}) ([]types.Business, error) {
	stmt := `Select * from "Business" WHERE organization_name=$1 OR rating=$3`
	ratingStr := query["rating"]
	rating, _ := strconv.Atoi(ratingStr.(string))
	rows, err := client.Conn.Query(stmt, query["organization_name"].(string), rating)

	if err != nil {
		return []types.Business{}, err
	}
	businesses := make([]types.Business, 0)
	defer rows.Close()
	for rows.Next() {
		var business types.Business

		err := rows.Scan(
			&business.OrganizationName, &business.Country, &business.AdministrativeArea, &business.Locality,
			&business.PostalCode,
			&business.Premise, &business.Thoroughfare, &business.ID, &business.Rating, &business.PricePerPound,
		)

		if err != nil {
			return businesses, err
		}

		businesses = append(businesses, business)
	}

	return businesses, err
}
