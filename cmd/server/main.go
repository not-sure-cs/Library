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

	godotenv.Load()

	portString := os.Getenv("PORT")

	if portString == "" {

		log.Fatal("PORT not Found")
	}

	fmt.Printf("Port: %s\n", portString)

	dbURL := os.Getenv("DB_URL")
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to Database")
	}

	db := database.New(conn)

	apiCfg := db

	mux := http.NewServeMux()

	//mux.HandleFunc("/Login", api.handleLogger())
	mux.HandleFunc("/status", api.HandleStatus(start))
	mux.HandleFunc("POST /book", api.HandleCreateBooks(apiCfg))
	mux.HandleFunc("GET /book/{name..}", api.HandleGetBooks(apiCfg))
	mux.HandleFunc("GET /book/{author..}", api.HandleListOfAuthorBooks(apiCfg))

	wrappedMux := api.JSONMiddleware(mux)

	srv := http.Server{
		Addr:    ":" + portString,
		Handler: wrappedMux,
	}

	fmt.Printf("Starting Server on Port: %s\n", portString)

	srv.ListenAndServe()

}
