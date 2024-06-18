package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/walrus811/blog-aggregator/internal/api"
)

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

	apiConfig := api.New(db)

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthz", api.HandlerReadiness)
	mux.HandleFunc("/v1/err", api.HandlerErr)

	mux.HandleFunc("GET /v1/users", apiConfig.MiddlewareAuth(apiConfig.HandlerGetUserByApiKey))
	mux.HandleFunc("POST /v1/users", apiConfig.HandlerCreateUser)

	mux.HandleFunc("GET /v1/feeds", apiConfig.HandlerGetFeeds)
	mux.HandleFunc("POST /v1/feeds", apiConfig.MiddlewareAuth(apiConfig.HandlerCreateFeed))

	mux.HandleFunc("GET /v1/feed_follows", apiConfig.MiddlewareAuth(apiConfig.HandlerGetFeedFollows))
	mux.HandleFunc("POST /v1/feed_follows", apiConfig.MiddlewareAuth(apiConfig.HandlerCreateFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiConfig.MiddlewareAuth(apiConfig.HandlerDeleteFeedFollow))

	http.ListenAndServe(":"+port, mux)
}
