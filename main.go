package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

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
