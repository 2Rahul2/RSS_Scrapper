package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/2Rahul2/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Hello world")
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found")
	}

	connection, err := sql.Open("postgres", dbURL) //returns a connection and err
	apiCfg := apiConfig{
		DB: database.New(connection),
	}

	if err != nil {
		log.Fatal("Could not connect to database")
	}

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
	V1router.Post("/users", apiCfg.handlerCreateUser)
	V1router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	V1router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))

	V1router.Get("/feeds", apiCfg.handlerGetFeeds)
	V1router.Post("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))

	V1router.Get("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))
	V1router.Delete("/feed-follows/{feedFollowId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", V1router)
	serve := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("server starting on port %v", portString)
	err = serve.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port : ", portString)
}
