package api

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const TEST_ENV_FILE = "../../.env.test"

func getTestApiConfig(fileName string) *ApiConfig {
	envLoadErr := godotenv.Load(fileName)
	if envLoadErr != nil {
		panic(envLoadErr)
	}

	dbURL := os.Getenv("DATABASE_URL")
	db, dbOpenERr := sql.Open("postgres", dbURL)

	if dbOpenERr != nil {
		panic(dbOpenERr)
	}

	return New(db)
}
