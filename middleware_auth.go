package main

import (
	"fmt"
	"github.com/MKefem/rssagg/internal/auth"
	"github.com/MKefem/rssagg/internal/database"
	"net/http"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r)
		if err != nil {
			respondWithError(
				w,
				403,
				fmt.Sprintf(
					"Auth error: %v",
					err,
				),
			)
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(
			r.Context(),
			apiKey,
		)
		if err != nil {
			respondWithError(
				w,
				400,
				fmt.Sprintf(
					"Error getting user: %v",
					err,
				),
			)
			return
		}
		handler(
			w,
			r,
			user,
		)
	}
}
