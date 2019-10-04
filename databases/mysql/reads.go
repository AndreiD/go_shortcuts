package database

import (
	"database/sql"
	"fmt"
)

// GetBalance .
func GetBalance(exchangeName, currency string) sql.NullFloat64 {
	var price sql.NullFloat64
	err := Db.Get(&price, "SELECT balance FROM balances WHERE exchange_name=? AND currency = ? LIMIT 1", exchangeName, currency)
	if err != nil {
		log.Error(err)
		return sql.NullFloat64{Float64: 0, Valid: false}
	}
	return price
}

// GetUserByEmail ...
func GetUserByEmail(email string) (models.User, error) {

	var user models.User
	err := Db.Get(&user, "SELECT * FROM users WHERE email=? LIMIT 1", email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return user, fmt.Errorf("invalid email or password")
		}
		return user, fmt.Errorf("can't query for the user id by email %s", err.Error())
	}
	return user, nil
}

func SomeValidation(user models.RegisterUser) error {

	// check if users's email exists
	var exists int
	err := Db.Get(&exists, "SELECT COUNT(id) FROM users WHERE email = ? LIMIT 1", user.Email)
	if err != nil {
		return fmt.Errorf("error %s", err.Error())
	}
	if exists > 0 {
		return fmt.Errorf("email already registered")
	}

	return nil
}

// GetAllUsersDB ..
func GetAllUsersDB(offset int, limit int) ([]models.User, error) {

	var users []models.User

	if limit > 2000 {
		return users, fmt.Errorf("limit is 2000 records per query")
	}
	if offset < 0 {
		return users, fmt.Errorf("invalid offset")
	}

	err := Db.Select(&users, "SELECT * FROM users LIMIT ?,?", offset, limit)
	if err != nil {
		return users, fmt.Errorf("can't query for users %s", err.Error())
	}

	return users, nil
}
