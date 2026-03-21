package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/knibirdgautam/library/internal/api"
	"github.com/knibirdgautam/library/internal/database"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	start := time.Now()

	godotenv.Load()

	portString := os.Getenv("PORT")

	if portString == "" {

		log.Fatal("PORT not Found")
	}

	fmt.Printf("Port: %s\n", portString)

	dbURL := os.Getenv("DB_URL")
	conn, err := database.NewConnection(dbURL)
	if err != nil {
		log.Fatal("Could not connect to DB:", err)
	}

	apiCfg := database.New(conn)

	mux := http.NewServeMux()

	//mux.HandleFunc("/Login", api.handleLogger())
	mux.HandleFunc("/status", api.HandleStatus(start))
	mux.HandleFunc("POST /book", api.HandleCreateBooks(apiCfg))
	mux.HandleFunc("GET /book/{id...}", api.HandleBookGet(apiCfg))
	mux.HandleFunc("PUT /book/{id}", api.HandleUpdateBooks(&apiCfg))
	mux.HandleFunc("DELETE /book/{id}", api.HandleDeleteBook(apiCfg))

	wrappedMux := api.JSONMiddleware(mux)

	srv := http.Server{
		Addr:    ":" + portString,
		Handler: wrappedMux,
	}

	fmt.Printf("Starting Server on Port: %s\n", portString)

	srv.ListenAndServe()

}
