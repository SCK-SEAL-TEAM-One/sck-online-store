package healthcheck

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

func GetUserNameFromDB(connection *sqlx.DB) (User, error) {
	user := User{}
	err := connection.Get(&user, "SELECT id,first_name,last_name FROM users WHERE id=1")
	if err != nil {
		fmt.Printf("Get user name from tearup get error : %s", err.Error())
		return User{}, err
	}
	return user, nil
}
