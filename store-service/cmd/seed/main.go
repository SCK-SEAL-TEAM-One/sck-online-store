package main

import (
	"log"
	"os"
	"store-service/internal/seed"
	"store-service/internal/user"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	dbConnection := "user:password@(localhost:3306)/store"
	if os.Getenv("DB_CONNECTION") != "" {
		dbConnection = os.Getenv("DB_CONNECTION")
	}

	connection, err := sqlx.Connect("mysql", dbConnection)
	if err != nil {
		log.Fatalln("cannot connect to database", err)
	}
	defer connection.Close()

	outputDir := "../shared"
	if os.Getenv("OUTPUT_DIR") != "" {
		outputDir = os.Getenv("OUTPUT_DIR")
	}

	seedUser := seed.SeedUserData{
		PasswordHelper: user.BcryptPasswordChecker{},
	}

	err = seedUser.GenerateUpdateUserDataCSV(outputDir, connection)
	if err != nil {
		log.Fatalf("Failed to generate CSV: %v", err)
	}

	log.Println("Seed CSV generation completed")
}
