package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/walrus811/blog-aggregator/internal/database"
	"github.com/walrus811/blog-aggregator/internal/utils"
)

const default_limit = 10

func (cfg *ApiConfig) HandlerGetPost(w http.ResponseWriter, r *http.Request, user database.User) {
	limit, limitParseErr := strconv.Atoi(r.URL.Query().Get("limit"))
	if limitParseErr != nil {
		limit = default_limit
	}
	posts, getPostsErr := cfg.GetPosts(r.Context(), limit)
	if getPostsErr != nil {
		log.Fatal(getPostsErr)
		utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, posts)
}
