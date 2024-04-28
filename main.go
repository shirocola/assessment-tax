package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	myhttp "github.com/shirocola/assessment-tax/pkg/http"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	router := mux.NewRouter()
	myhttp.RegisterRoutes(router)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("Shutting down the server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Error during shutdown: %v", err)
		}
	}()

	protectedRoutes := router.PathPrefix("/admin").Subrouter()
	protectedRoutes.Use(myhttp.BasicAuthMiddleware)
	protectedRoutes.HandleFunc("/deductions/personal", myhttp.SetPersonalDeductionHandler).Methods("POST")
	protectedRoutes.HandleFunc("/deductions/k-receipt", myhttp.SetKReceiptDeductionHandler).Methods("POST")

	log.Printf("Starting server on port %s", os.Getenv("PORT"))
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server stopped unexpectedly: %v", err)
	}

	log.Println("Server stopped")
}
