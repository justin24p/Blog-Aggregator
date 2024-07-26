package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/justin24p/rssAggregator/internal/database"
	// _ "github.com/lib/pg"
)

// hold database connection
type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("hello world")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port isn't found")
	}

	// import database connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found")
	}
	// conneciton to db
	conn, err := sql.Open("postgres",dbURL)
	if err != nil{
		log.Fatal("Cant connect to database")
	}

	apiCfg := apiConfig{
		DB: database.New(conn), 
	}


	router := chi.NewRouter()

	// allows request from broswer 
	// just sends extra headers allowing user to have more control
	router.Use(cors.Handler(cors.Options{
	  	AllowedOrigins:   []string{"http://*"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, 
	}))	

	// /v1/ready route  
	// gives us two handlers 
	// scopes for get
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz",handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users",apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feeds",apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds",apiCfg.handlerGetFeeds)
	
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))

	router.Mount("/v1",v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	// listenandserver blocks 
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}