package main

import (
	"database/sql"
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

	err := godotenv.Load()

	if err != nil {
        log.Println("Error loading .env file, proceeding without it")
    }

	portString := os.Getenv("PORT")

	if portString == "" {

		log.Fatal("PORT not found")
	}

	fmt.Printf("Port: %s\n", portString)

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
        log.Fatal("DB_URL not found")
    }

	fmt.Printf("DB_URL: %s\n", dbURL)

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {

		log.Fatal("Cant connect to Database")
	}

	    if err := conn.Ping(); err != nil {
        log.Fatalf("Cannot reach Database: %v", err)
    }

	 fmt.Println("Successfully connected to the database")

	db := database.New(conn)

	apiCfg := db

	mux := http.NewServeMux()

	//mux.HandleFunc("/Login", api.handleLogger())
	mux.HandleFunc("/status", api.HandleStatus(start))
	mux.HandleFunc("/book", api.HandleCreateBooks(apiCfg))
	mux.HandleFunc("/book/{id}}", api.HandleGetBooks(apiCfg))
	mux.HandleFunc("/book/{id}", api.HandleListOfAuthorBooks(apiCfg))

	wrappedMux := api.JSONMiddleware(mux)

	srv := http.Server{
		Addr:    ":" + portString,
		Handler: wrappedMux,
	}

	fmt.Printf("Starting Server on Port: %s\n", portString)

	srv.ListenAndServe()

}
