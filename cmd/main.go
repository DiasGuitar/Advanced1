package main

import (
	"assignment1/internal/data"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pressly/goose"
	"log"
	"net/http"
	"time"

	"database/sql"
	// Import PostgreSQL driver
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

// My connection string:
// "host=localhost port=5432 user=postgres password=123456 dbname=d.mukhamedinDB sslmode=disable"

func main() {
	// DB connection
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

	// goose -dir migrations postgres "host=localhost port=5432 user=postgres password=123456 dbname=d.mukhamedinDB sslmode=disable" up
	// goose -dir migrations postgres "host=localhost port=5432 user=postgres password=123456 dbname=d.mukhamedinDB sslmode=disable" down

	r := mux.NewRouter()
	r.HandleFunc("/module/{id}", getModuleInfo).Methods("GET")
	r.HandleFunc("/module", createModuleInfo).Methods("POST")
	r.HandleFunc("/module/{id}", updateModuleInfo).Methods("PUT")
	r.HandleFunc("/module/{id}", deleteModuleInfo).Methods("DELETE")

	log.Println("Server is listening on :8080")
	http.ListenAndServe(":8080", r)
}

func createModuleInfo(w http.ResponseWriter, r *http.Request) {
	var module data.ModuleInfo
	err := json.NewDecoder(r.Body).Decode(&module)
	if err != nil {
		http.Error(w, "Error during module reading", http.StatusBadRequest)
		return
	}

	module.CreatedAt = time.Now()
	module.UpdatedAt = time.Now()

	// Execute SQL statement to insert data
	_, err = db.Exec("INSERT INTO module_info (module_name, module_duration, exam_type, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		module.ModuleName, module.ModuleDuration, module.ExamType, module.Version, module.CreatedAt, module.UpdatedAt)
	if err != nil {
		http.Error(w, "Error during module creation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(module)
}

func getModuleInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var module data.ModuleInfo

	// Execute SQL query to fetch data
	row := db.QueryRow("SELECT * FROM module_info WHERE id = $1", id)
	err := row.Scan(&module.ID, &module.ModuleName, &module.ModuleDuration, &module.ExamType, &module.Version, &module.CreatedAt, &module.UpdatedAt)
	if err != nil {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(module)
}

func updateModuleInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var module data.ModuleInfo

	// Execute SQL query to fetch data
	row := db.QueryRow("SELECT * FROM module_info WHERE id = $1", id)
	err := row.Scan(&module.ID, &module.ModuleName, &module.ModuleDuration, &module.ExamType, &module.Version, &module.CreatedAt, &module.UpdatedAt)
	if err != nil {
		http.Error(w, "Module not found", http.StatusNotFound)
		return
	}

	var updatedModule data.ModuleInfo
	err = json.NewDecoder(r.Body).Decode(&updatedModule)
	if err != nil {
		http.Error(w, "Error retrieving module", http.StatusBadRequest)
		return
	}

	// Update the module info
	module.ModuleName = updatedModule.ModuleName
	module.ModuleDuration = updatedModule.ModuleDuration
	module.ExamType = updatedModule.ExamType
	module.Version = updatedModule.Version
	module.UpdatedAt = time.Now()

	// Execute SQL statement to update data
	_, err = db.Exec("UPDATE module_info SET module_name=$1, module_duration=$2, exam_type=$3, version=$4, updated_at=$5 WHERE id=$6",
		module.ModuleName, module.ModuleDuration, module.ExamType, module.Version, module.UpdatedAt, id)
	if err != nil {
		http.Error(w, "Error updating module", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(module)
}

func deleteModuleInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// Execute SQL statement to delete data
	_, err := db.Exec("DELETE FROM module_info WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Error deleting module", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
