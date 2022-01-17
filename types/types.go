package types

import "time"

type Addressable struct {
	Country            string `json:"country"`
	AdministrativeArea string `json:"administrative_area"`
	Thoroughfare       string `json:"thoroughfare"`
	Premise            string `json:"premise"`
	Locality           string `json:"locality"`
	PostalCode         string `json:"postal_code"`
}

type Business struct {
	Addressable
	ID               int     `json:"id"`
	OrganizationName string  `json:"organization_name"`
	Rating           int     `json:"rating"`
	PricePerPound    float64 `json:"price_per_pound"`
}

type User struct {
	Addressable
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Order struct {
	ID                 string  `json:"id"`
	UserID             int     `json:"user_id"`
	BusinessID         int     `json:"business_id"`
	UseBleach          bool    `json:"use_bleach"`
	Status             string  `json:"status"`
	PreferredDetergent string  `json:"preferred_detergent"`
	PreferredSoftener  string  `json:"preferred_softener"`
	WeightLBS          float64 `json:"weight_lbs"`
	PickupCharge       float64 `json:"pickup_charge"`
	DropCharge         float64 `json:"drop_charge"`
	ServiceCharge      float64 `json:"service_charge"`
}

type Session struct {
	ID                    int       `json:"id"`
	UserID                int       `json:"user_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiration time.Time `json:"access_token_expiration"`
}
