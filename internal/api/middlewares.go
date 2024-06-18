package api

import (
	"net/http"

	"github.com/walrus811/blog-aggregator/internal/database"
	"github.com/walrus811/blog-aggregator/internal/utils"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, authErr := cfg.CheckAuth(r.Context(), r.Header)
		if authErr != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		handler(w, r, user)
	}
}
