package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello world")
	godotenv.Load(".env")
	portString := os.Getenv("PORT")

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	V1router := chi.NewRouter()
	V1router.Get("/ready", handlerReadiness)
	V1router.Get("/err", handlerErr)

	router.Mount("/v1", V1router)
	serve := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("server starting on port %v", portString)
	err := serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	if portString == "" {
		log.Fatal("Port is not found")
	}

	fmt.Println("Port : ", portString)
}
