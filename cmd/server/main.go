package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/knibirdgautam/library/internal/api"
	"github.com/knibirdgautam/library/internal/database"
	"github.com/knibirdgautam/library/internal/storage"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func init() {
	gob.Register(uuid.UUID{})
	gob.Register(sql.NullString{})
}

func main() {
	start := time.Now()

	err := godotenv.Load()
	if err != nil {
		log.Println("Notice: .env file not found or couldn't be loaded, using process environment variables")
	}

	sessionKey := os.Getenv("KEY")
	if sessionKey == "" {
		log.Fatal("KEY not found in environment")
	}
	store := sessions.NewCookieStore([]byte(sessionKey))

	isProd := os.Getenv("ENV") == "production"
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
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

	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Cant connect to Database")
	}
	defer conn.Close()

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err := conn.Ping(); err != nil {
		log.Fatalf("Cannot reach Database: %v", err)
	}

	fmt.Println("Successfully connected to the database")

	db := database.New(conn)

	apiCfg := db

	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")

	var config = storage.Secret{
		Key:     secretKey,
		Access:  accessKey,
		Account: accountID,
		Bucket:  bucketName,
	}

	R2Cclient, err := storage.NewR2Client(accountID, accessKey, secretKey)
	if err != nil {
		log.Fatalf("Cannot Connect to Client: %v", err)
	}

	var client = storage.R2Store{
		Client: R2Cclient,
	}

	mux := http.NewServeMux()

	requireModOrAdmin := api.RequireRoles(store, database.UserRoleModerator, database.UserRoleAdmin)

	mux.HandleFunc("GET /status", api.HandleStatus(start))
	mux.HandleFunc("POST /user/signup", api.HandleSignUp(apiCfg))
	mux.HandleFunc("POST /user/login", api.HandleLogging(apiCfg, store))

	mux.Handle("POST /book", api.AuthedMiddleware(api.HandleCreateBooks(apiCfg, client, config), store))
	mux.Handle("GET /book/{id}", api.AuthedMiddleware(api.HandleGetBooks(apiCfg, client, config), store))
	mux.Handle("PUT /book/{id}", api.AuthedMiddleware(requireModOrAdmin(api.HandleUpdateBook(apiCfg)), store))
	mux.Handle("DELETE /book/{id}", api.AuthedMiddleware(requireModOrAdmin(api.HandleDeleteBook(apiCfg, client, config)), store))

	wrappedMux := api.JSONMiddleware(mux)

	srv := http.Server{
		Addr:              ":" + portString,
		Handler:           wrappedMux,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      10 * time.Minute,
		IdleTimeout:       3 * time.Minute,
		ReadHeaderTimeout: 5 * time.Minute,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		fmt.Printf("Starting Server on Port: %s\n", portString)
		serverErrors <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("[Server] Server error: %v", err)
		}
	case <-shutdownCtx.Done():
		log.Println("[Server] Shutdown signal received, gracefully terminating...")

		stopCtx, cancelStop := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelStop()

		if err := srv.Shutdown(stopCtx); err != nil {
			log.Printf("[Server] Forced server shutdown error: %v\n", err)
		} else {
			log.Println("[Server] HTTP server stopped cleanly")
		}
	}

	log.Println("[Server] Server stopped")
}
