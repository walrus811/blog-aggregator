package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/walrus811/blog-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	envLoadErr := godotenv.Load()

	if envLoadErr != nil {
		panic(envLoadErr)
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")

	db, dbOpenERr := sql.Open("postgres", dbURL)

	if dbOpenERr != nil {
		panic(dbOpenERr)
	}

	dbQueries := database.New(db)

	apiConfig := &apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthz", handlerReadiness)
	mux.HandleFunc("/v1/err", handlerErr)

	mux.HandleFunc("GET /v1/users", apiConfig.middlewareAuth(handlerGetUserByApiKey))
	mux.HandleFunc("POST /v1/users", apiConfig.handlerCreateUser)

	mux.HandleFunc("GET /v1/feeds",apiConfig.handlerGetFeeds)
	mux.HandleFunc("POST /v1/feeds", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))

	http.ListenAndServe(":"+port, mux)
}
