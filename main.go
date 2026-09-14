package main

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/montruh-afk/chirpy/internal"
	"github.com/montruh-afk/chirpy/internal/database"
	"database/sql"
)




func main () {
	godotenv.Load()
	dbURL := os.Getenv("DBURL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Something went wrong while attemptiong to access our records: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)


	handler := http.NewServeMux()

	s := &http.Server{
		Addr: os.Getenv("PORT"),
		Handler: handler,
	}
	cfg := &internal.ApiConfig{
		FileServerHits: atomic.Int32{},
		Db: dbQueries,
	}

	startUp(handler, cfg)

	
	log.Printf("Serving on port %s...\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}