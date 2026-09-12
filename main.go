package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/montruh-afk/chirpy/internal"
)




func main () {
	handler := http.NewServeMux()

	s := &http.Server{
		Addr: ":8080",
		Handler: handler,
	}
	cfg := &internal.ApiConfig{
		FileServerHits: atomic.Int32{},
	}

	startUp(handler, cfg)

	
	log.Printf("Serving on port %s...\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}