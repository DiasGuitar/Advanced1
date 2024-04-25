package handlers

import (
	"assignment1/internal/data"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateUserInfoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Your implementation here
		var userInfo data.UserInfo
		err := json.NewDecoder(r.Body).Decode(&userInfo)
		if err != nil {
			http.Error(w, "Error during user reading", http.StatusBadRequest)
			return
		}

		// Execute SQL statement to insert data
		_, err = db.Exec("INSERT INTO user_info (name, surname, email, password_hash, role, activated, version) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			userInfo.Name, userInfo.Surname, userInfo.Email, userInfo.PasswordHash, userInfo.Role, userInfo.Activated, userInfo.Version)
		if err != nil {
			http.Error(w, "Error during user creation", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(userInfo)
	}
}

func GetUserInfoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Your implementation here
		params := mux.Vars(r)
		id := params["id"]

		var userInfo data.UserInfo

		// Execute SQL query to fetch data
		row := db.QueryRow("SELECT * FROM user_info WHERE id = $1", id)
		err := row.Scan(&userInfo.ID, &userInfo.CreatedAt, &userInfo.UpdatedAt, &userInfo.Name, &userInfo.Surname, &userInfo.Email, &userInfo.PasswordHash, &userInfo.Role, &userInfo.Activated, &userInfo.Version)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(userInfo)
	}
}

func EditUserInfoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Your implementation here
		params := mux.Vars(r)
		id := params["id"]

		var updatedUserInfo data.UserInfo
		err := json.NewDecoder(r.Body).Decode(&updatedUserInfo)
		if err != nil {
			http.Error(w, "Error retrieving user", http.StatusBadRequest)
			return
		}

		// Update the user info
		// Execute SQL statement to update data
		_, err = db.Exec("UPDATE user_info SET name=$1, surname=$2, email=$3, password_hash=$4, role=$5, activated=$6, version=$7 WHERE id=$8",
			updatedUserInfo.Name, updatedUserInfo.Surname, updatedUserInfo.Email, updatedUserInfo.PasswordHash, updatedUserInfo.Role, updatedUserInfo.Activated, updatedUserInfo.Version, id)
		if err != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(updatedUserInfo)
	}
}

func DeleteUserInfoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Your implementation here
		params := mux.Vars(r)
		id := params["id"]

		// Execute SQL statement to delete data
		_, err := db.Exec("DELETE FROM user_info WHERE id=$1", id)
		if err != nil {
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
