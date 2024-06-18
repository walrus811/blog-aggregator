package api

import (
	"database/sql"

	"github.com/walrus811/blog-aggregator/internal/database"
)


type ApiConfig struct {
	DB      *database.Queries
	DBInner *sql.DB
}

func New(db *sql.DB) *ApiConfig {
	return &ApiConfig{
		DB:      database.New(db),
		DBInner: db,
	}
}