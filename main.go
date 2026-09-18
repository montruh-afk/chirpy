package main

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/montruh-afk/chirpy/internal"
	"github.com/montruh-afk/chirpy/internal/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DBURL")
	if dbURL == "" {
		log.Fatal("DBURL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Something went wrong while attemptiong to access our records: %s", err)
		os.Exit(1)
	}

	defer db.Close()


	//ping first to avoid panic after server has started up
	if err := db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	dbQueries := database.New(db)

	handler := http.NewServeMux()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	s := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	cfg := &internal.ApiConfig{
		FileServerHits: atomic.Int32{},
		Db:             dbQueries,
		Platform:       os.Getenv("PLATFORM"),
		TknScrt: os.Getenv("HASHTEXT"),
		Exp: os.Getenv("EXP"),
	}

	if len(cfg.Exp) < 1 {
		cfg.Exp = "3600s"
	}

	startUp(handler, cfg)

	log.Printf("Serving on port: %s...\n", port)
	log.Fatal(s.ListenAndServe())
}
