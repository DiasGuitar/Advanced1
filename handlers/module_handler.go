package handlers

import (
	"assignment1/internal/data"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func CreateModuleInfo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var module data.ModuleInfo
		err := json.NewDecoder(r.Body).Decode(&module)
		if err != nil {
			http.Error(w, "Error during module reading", http.StatusBadRequest)
			return
		}

		module.CreatedAt = time.Now()
		module.UpdatedAt = time.Now()

		_, err = db.Exec("INSERT INTO module_info (module_name, module_duration, exam_type, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
			module.ModuleName, module.ModuleDuration, module.ExamType, module.Version, module.CreatedAt, module.UpdatedAt)
		if err != nil {
			http.Error(w, "Error during module creation", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(module)
	}
}

func GetModuleInfo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var module data.ModuleInfo

		row := db.QueryRow("SELECT * FROM module_info WHERE id = $1", id)
		err := row.Scan(&module.ID, &module.ModuleName, &module.ModuleDuration, &module.ExamType, &module.Version, &module.CreatedAt, &module.UpdatedAt)
		if err != nil {
			http.Error(w, "Module not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(module)
	}
}

func UpdateModuleInfo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		var module data.ModuleInfo

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

		module.ModuleName = updatedModule.ModuleName
		module.ModuleDuration = updatedModule.ModuleDuration
		module.ExamType = updatedModule.ExamType
		module.Version = updatedModule.Version
		module.UpdatedAt = time.Now()

		_, err = db.Exec("UPDATE module_info SET module_name=$1, module_duration=$2, exam_type=$3, version=$4, updated_at=$5 WHERE id=$6",
			module.ModuleName, module.ModuleDuration, module.ExamType, module.Version, module.UpdatedAt, id)
		if err != nil {
			http.Error(w, "Error updating module", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(module)
	}
}

func DeleteModuleInfo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		_, err := db.Exec("DELETE FROM module_info WHERE id=$1", id)
		if err != nil {
			http.Error(w, "Error deleting module", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
