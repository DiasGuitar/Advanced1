package handlers

import (
	"assignment1/internal/data"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	dsn := "host=localhost user=postgres password=123456 dbname=test_db port=5432 sslmode=disable"
	var err error
	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	m.Run()
}

func TestCreateModuleInfo(t *testing.T) {
	payload := data.ModuleInfo{
		ModuleName:     "Test Module",
		ModuleDuration: 100,
		ExamType:       "Test",
		Version:        "1.0",
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/module", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := CreateModuleInfo(testDB)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var module data.ModuleInfo
	json.NewDecoder(rr.Body).Decode(&module)

	assert.Equal(t, payload.ModuleName, module.ModuleName)
}

func TestGetModuleInfo(t *testing.T) {
	module := data.ModuleInfo{
		ModuleName:     "Test Module",
		ModuleDuration: 100,
		ExamType:       "Test",
		Version:        "1.0",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := testDB.Exec("INSERT INTO module_info (module_name, module_duration, exam_type, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		module.ModuleName, module.ModuleDuration, module.ExamType, module.Version, module.CreatedAt, module.UpdatedAt)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/module/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := GetModuleInfo(testDB)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedModule data.ModuleInfo
	json.NewDecoder(rr.Body).Decode(&retrievedModule)

	assert.Equal(t, module.ModuleName, retrievedModule.ModuleName)
}

func TestUpdateModuleInfo(t *testing.T) {
	module := data.ModuleInfo{
		ModuleName:     "Test Module",
		ModuleDuration: 100,
		ExamType:       "Test",
		Version:        "1.0",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := testDB.Exec("INSERT INTO module_info (module_name, module_duration, exam_type, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		module.ModuleName, module.ModuleDuration, module.ExamType, module.Version, module.CreatedAt, module.UpdatedAt)
	assert.NoError(t, err)

	updatedModule := data.ModuleInfo{
		ModuleName:     "Updated Module",
		ModuleDuration: 200,
		ExamType:       "Final",
		Version:        "2.0",
	}

	body, _ := json.Marshal(updatedModule)
	req, err := http.NewRequest("PUT", "/module/1", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := UpdateModuleInfo(testDB)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedModule data.ModuleInfo
	json.NewDecoder(rr.Body).Decode(&retrievedModule)

	assert.Equal(t, updatedModule.ModuleName, retrievedModule.ModuleName)
}

func TestDeleteModuleInfo(t *testing.T) {
	module := data.ModuleInfo{
		ModuleName:     "Test Module",
		ModuleDuration: 100,
		ExamType:       "Test",
		Version:        "1.0",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := testDB.Exec("INSERT INTO module_info (module_name, module_duration, exam_type, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		module.ModuleName, module.ModuleDuration, module.ExamType, module.Version, module.CreatedAt, module.UpdatedAt)
	assert.NoError(t, err)

	req, err := http.NewRequest("DELETE", "/module/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := DeleteModuleInfo(testDB)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestCreateUserInfoHandler(t *testing.T) {
	payload := data.UserInfo{
		Name:         "John",
		Surname:      "Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "password",
		Role:         "user",
		Activated:    true,
		Version:      1,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := CreateUserInfoHandler(testDB)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var user data.UserInfo
	json.NewDecoder(rr.Body).Decode(&user)

	assert.Equal(t, payload.Name, user.Name)
}

func TestGetUserInfoHandler(t *testing.T) {
	user := data.UserInfo{
		Name:         "John",
		Surname:      "Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "password",
		Role:         "user",
		Activated:    true,
		Version:      1,
	}

	_, err := testDB.Exec("INSERT INTO user_info (name, surname, email, password_hash, role, activated, version) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.Name, user.Surname, user.Email, user.PasswordHash, user.Role, user.Activated, user.Version)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/user/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := GetUserInfoHandler(testDB)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedUser data.UserInfo
	json.NewDecoder(rr.Body).Decode(&retrievedUser)

	assert.Equal(t, user.Name, retrievedUser.Name)
}

func TestEditUserInfoHandler(t *testing.T) {
	user := data.UserInfo{
		Name:         "John",
		Surname:      "Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "password",
		Role:         "user",
		Activated:    true,
		Version:      1,
	}

	_, err := testDB.Exec("INSERT INTO user_info (name, surname, email, password_hash, role, activated, version) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.Name, user.Surname, user.Email, user.PasswordHash, user.Role, user.Activated, user.Version)
	assert.NoError(t, err)

	updatedUser := data.UserInfo{
		Name:         "Jane",
		Surname:      "Doe",
		Email:        "jane.doe@example.com",
		PasswordHash: "newpassword",
		Role:         "admin",
		Activated:    false,
		Version:      2,
	}

	body, _ := json.Marshal(updatedUser)
	req, err := http.NewRequest("PUT", "/user/1", bytes.NewBuffer(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := EditUserInfoHandler(testDB)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedUser data.UserInfo
	json.NewDecoder(rr.Body).Decode(&retrievedUser)

	assert.Equal(t, updatedUser.Name, retrievedUser.Name)
}

func TestDeleteUserInfoHandler(t *testing.T) {
	user := data.UserInfo{
		Name:         "John",
		Surname:      "Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "password",
		Role:         "user",
		Activated:    true,
		Version:      1,
	}

	_, err := testDB.Exec("INSERT INTO user_info (name, surname, email, password_hash, role, activated, version) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		user.Name, user.Surname, user.Email, user.PasswordHash, user.Role, user.Activated, user.Version)
	assert.NoError(t, err)

	req, err := http.NewRequest("DELETE", "/user/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := DeleteUserInfoHandler(testDB)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}
