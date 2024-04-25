package main

import (
	"assignment1/handlers"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gorilla/mux"
	"github.com/pressly/goose"

	"database/sql"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "d.mukhamedinDB"
)

var db *sql.DB

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, dbname, port)
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	// Check if the connection to the database can be established
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection establishing error: %v", err)
	}

	// Applying migrations
	err = goose.Up(db, "./migrations")
	if err != nil {
		log.Fatalf("Migration applying error: %v", err)
	}

	// Sending welcome email
	err = sendWelcomeEmail("brainwin0@gmail.com")
	if err != nil {
		log.Println("Error sending welcome email:", err)
	}

	fmt.Println("Welcome email sent successfully")

	r := mux.NewRouter()

	// Module handlers
	moduleRouter := r.PathPrefix("/module").Subrouter()
	moduleRouter.HandleFunc("/{id}", handlers.GetModuleInfo(db)).Methods("GET")
	moduleRouter.HandleFunc("/", handlers.CreateModuleInfo(db)).Methods("POST")
	moduleRouter.HandleFunc("/{id}", handlers.UpdateModuleInfo(db)).Methods("PUT")
	moduleRouter.HandleFunc("/{id}", handlers.DeleteModuleInfo(db)).Methods("DELETE")

	// User handlers
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/", handlers.CreateUserInfoHandler(db)).Methods("POST")
	userRouter.HandleFunc("/{id}", handlers.GetUserInfoHandler(db)).Methods("GET")
	userRouter.HandleFunc("/{id}", handlers.EditUserInfoHandler(db)).Methods("PUT")
	userRouter.HandleFunc("/{id}", handlers.DeleteUserInfoHandler(db)).Methods("DELETE")

	log.Println("Server is listening on :8080")
	http.ListenAndServe(":8080", r)
}

func sendWelcomeEmail(email string) error {

	from := os.Getenv("SENDER_EMAIL")
	if from == "" {
		return fmt.Errorf("SENDER_EMAIL environment variable is not set")
	}
	fmt.Println(os.Getenv("SENDER_EMAIL"))

	password := os.Getenv("SENDER_PASSWORD")
	if password == "" {
		return fmt.Errorf("SENDER_PASSWORD environment variable is not set")
	}

	to := email

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP_HOST or SMTP_PORT environment variables are not set")
	}

	message := []byte(fmt.Sprintf("Welcome!"))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return err
	}

	return nil
}
