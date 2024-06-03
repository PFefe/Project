package main

import (
	"database/sql"
	"fmt"
	"github.com/MKefem/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open(
		"postgres",
		dbURL,
	)
	if err != nil {
		log.Fatal(
			"Failed to connect to the database: %v",
			err,
		)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()
	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"htts//*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"*"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300,
			},
		),
	)

	v1Router := chi.NewRouter()
	v1Router.Get(
		"/healthz",
		handlerReadiness,
	)
	v1Router.Get(
		"/err",
		handlerErr,
	)
	v1Router.Post(
		"/users",
		apiCfg.handlerCreateUser,
	)
	v1Router.Get(
		"/users",
		apiCfg.middlewareAuth(apiCfg.handlerGetUser),
	)
	v1Router.Post(
		"/feeds",
		apiCfg.middlewareAuth(apiCfg.handlerCreateFeed),
	)
	v1Router.Get(
		"/feeds",
		apiCfg.handlerGetFeeds,
	)
	v1Router.Post(
		"/feed_follows",
		apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollows),
	)
	v1Router.Get(
		"/feed_follows",
		apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows),
	)
	v1Router.Delete(
		"/feed_follows/{feedFollowID}",
		apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow),
	)

	router.Mount(
		"/v1",
		v1Router,
	)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf(
		"Server is running on port %v",
		portString,
	)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(
		"Port:",
		portString,
	)
}
