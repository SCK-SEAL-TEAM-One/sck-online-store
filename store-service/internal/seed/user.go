package seed

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"store-service/internal/user"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID int `json:"id" db:"id"`
}

func GenerateUpdateUserDataCSV(outputDir string, db *sqlx.DB) error {
	filePath := filepath.Join(outputDir, "001-users-with-username-password.csv")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"id", "username", "password"})
	if err != nil {
		return err
	}

	var users []User
	query := `SELECT id FROM users`
	err = db.Select(&users, query)
	if err != nil {
		return err
	}

	log.Printf("Users found in DB: %d", len(users))

	defaultPassword := "P@ssw0rd"
	for _, u := range users {
		hashed, err := user.HashPassword(defaultPassword)
		if err != nil {
			return err
		}

		username := fmt.Sprintf("user_%d", u.ID)
		log.Printf("username: %s", username)

		record := []string{
			fmt.Sprintf("%d", u.ID),
			username,
			hashed,
		}

		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	log.Println("CSV generated: 001-users-with-username-password.csv")
	return nil
}
