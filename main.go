package main

import (
	"log"
	"net/http"

	"github.com/montruh-afk/chirpy/internal"
)




func main () {
	handler := http.NewServeMux()

	s := &http.Server{
		Addr: ":8080",
		Handler: handler,
	}

	root := http.Dir(".")
	reqPath := http.StripPrefix("/app", http.FileServer(root))
	handler.Handle("/app/", reqPath)
	handler.HandleFunc("/healthz", internal.ReadinessEndpoint)
	log.Printf("Serving on port %s\n", s.Addr)
	log.Fatal(s.ListenAndServe())
}