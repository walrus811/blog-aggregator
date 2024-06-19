package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/walrus811/blog-aggregator/internal/api"
	"github.com/walrus811/blog-aggregator/internal/workers"
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
	log.Println("Connected to database")
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

	mux.HandleFunc("GET /v1/posts", apiConfig.MiddlewareAuth(apiConfig.HandlerGetPost))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	const collectionConcurrency = 10
	const collectionInterval = 10 * time.Second
	go workers.StartScrape(apiConfig.DB, collectionConcurrency, collectionInterval)

	log.Println("Server started on port", port)
	log.Fatal(srv.ListenAndServe())
}
