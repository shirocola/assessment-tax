package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	myhttp "github.com/shirocola/assessment-tax/pkg/http"
)

func main() {
	router := mux.NewRouter()

	myhttp.RegisterRoutes(router)

	port := "8080"
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
