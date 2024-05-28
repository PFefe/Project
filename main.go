package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found in the environment")
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

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(
		"Port:",
		portString,
	)
}
