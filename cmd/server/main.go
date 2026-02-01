package main

import (
	"fmt"
	"lib/internal/api"
	"lib/internal/database"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	start := time.Now()

	godotenv.Load()

	portString := os.Getenv("PORT")

	if portString == "" {

		log.Fatal("PORT not Found")
	}

	fmt.Printf("Port: %s\n", portString)

	mux := http.NewServeMux()

	mux.HandleFunc("/status", api.HandleStatus(start))
	mux.HandleFunc("POST /book", api.HandleCreateBooks(&database.BookDB))
	mux.HandleFunc("GET /book/{id...}", api.HandleBookGet(&database.BookDB))
	mux.HandleFunc("PUT /book/{id}", api.HandleUpdateBooks(&database.BookDB))
	mux.HandleFunc("DELETE /book/{id}", api.HandleDeleteBook(&database.BookDB))

	wrappedMux := api.JSONMiddleware(mux)

	srv := http.Server{
		Addr:    ":" + portString,
		Handler: wrappedMux,
	}

	fmt.Printf("Starting Server on Port: %s\n", portString)

	srv.ListenAndServe()

}
