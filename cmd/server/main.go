package main

import (
	"log"
	"net/http" // Standard library's HTTP package

	myhttp "github.com/shirocola/assessment-tax/pkg/http" // Your custom HTTP package
)

func main() {
	http.HandleFunc("/tax/calculations", myhttp.TaxCalculationHandler) // Use the alias here
	port := "8080"                                                     // Example, typically you'd get this from environment variables

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
