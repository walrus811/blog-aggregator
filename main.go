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
	mux := serverMux(*apiConfig)

	http.ListenAndServe(":"+port, mux)
}

func serverMux(cfg api.ApiConfig) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthz", api.HandlerReadiness)
	mux.HandleFunc("/v1/err", api.HandlerErr)

	mux.HandleFunc("GET /v1/users", cfg.MiddlewareAuth(cfg.HandlerGetUserByApiKey))
	mux.HandleFunc("POST /v1/users", cfg.HandlerCreateUser)

	mux.HandleFunc("GET /v1/feeds", cfg.HandlerGetFeeds)
	mux.HandleFunc("POST /v1/feeds", cfg.MiddlewareAuth(cfg.HandlerCreateFeed))

	mux.HandleFunc("GET /v1/feed_follows", cfg.MiddlewareAuth(cfg.HandlerGetFeedFollows))
	mux.HandleFunc("POST /v1/feed_follows", cfg.MiddlewareAuth(cfg.HandlerCreateFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.MiddlewareAuth(cfg.HandlerDeleteFeedFollow))

	return mux
}
