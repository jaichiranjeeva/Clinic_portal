// integration_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"portal/config"
	"portal/controllers"
	"portal/middleware"
	"portal/models"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testRouter *gin.Engine
var doctorToken, receptionistToken string

func initTestDB() {
	dsn := "host=localhost user=postgres password=123 dbname=hospital_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Database connection failed")
	}
	db.Migrator().DropTable(&models.User{}, &models.Patient{})
	db.AutoMigrate(&models.User{}, &models.Patient{})
	config.DB = db
}

func setupRouter() {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.Group("/receptionist").Use(middleware.AuthMiddleware("receptionist")).POST("/patients", controllers.CreatePatient)
	r.Group("/receptionist").Use(middleware.AuthMiddleware("receptionist")).DELETE("/patients/:id", controllers.DeletePatient)
	doctor := r.Group("/doctor")
	doctor.Use(middleware.AuthMiddleware("doctor"))
	{
		doctor.GET("/patients", controllers.GetPatients)
		doctor.GET("/patients/:id", controllers.FetchPatient)
		doctor.PUT("/patients/:id", controllers.UpdatePatient)
		doctor.DELETE("/patients/:id", controllers.DeletePatient)
	}
	testRouter = r
}

func TestRolesAndPatientWorkflow(t *testing.T) {
	initTestDB()
	setupRouter()

	registerAndLogin(t, map[string]string{"username": "doc", "password": "123", "role": "doctor"}, &doctorToken)
	registerAndLogin(t, map[string]string{"username": "recep", "password": "123", "role": "receptionist"}, &receptionistToken)

	// Receptionist creates patient
	patient := map[string]interface{}{
		"name":   "Asha",
		"age":    25,
		"gender": "Female",
		"note":   "Routine",
	}
	body, _ := json.Marshal(patient)
	req := httptest.NewRequest("POST", "/receptionist/patients", bytes.NewBuffer(body))
	req.Header.Set("Authorization", receptionistToken)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Receptionist failed to create patient, got %d", resp.Code)
	}

	// Doctor fetches patients
	req = httptest.NewRequest("GET", "/doctor/patients", nil)
	req.Header.Set("Authorization", doctorToken)
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Doctor failed to get patients, got %d", resp.Code)
	}

	// Doctor updates patient with valid age
	update := map[string]interface{}{
		"name":   "Updated",
		"age":    40,
		"gender": "Female",
		"note":   "Updated Note",
	}
	body, _ = json.Marshal(update)
	req = httptest.NewRequest("PUT", "/doctor/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", doctorToken)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Update failed with valid age, got %d", resp.Code)
	}

	// Doctor update with invalid age
	update["age"] = 0
	body, _ = json.Marshal(update)
	req = httptest.NewRequest("PUT", "/doctor/patients/1", bytes.NewBuffer(body))
	req.Header.Set("Authorization", doctorToken)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	if resp.Code == http.StatusOK {
		t.Error("Update should fail with invalid age")
	}

	// Receptionist deletes patient
	req = httptest.NewRequest("DELETE", "/receptionist/patients/1", nil)
	req.Header.Set("Authorization", receptionistToken)
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Errorf("Receptionist failed to delete patient, got %d", resp.Code)
	}

	// Unauthorized: no token
	req = httptest.NewRequest("GET", "/doctor/patients", nil)
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for missing token, got %d", resp.Code)
	}

	// Invalid login
	invalid := map[string]string{"username": "wrong", "password": "fail"}
	body, _ = json.Marshal(invalid)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for invalid login, got %d", resp.Code)
	}
}

func registerAndLogin(t *testing.T, user map[string]string, token *string) {
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	loginPayload := map[string]string{"username": user["username"], "password": user["password"]}
	body, _ = json.Marshal(loginPayload)
	req = httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	var loginResp map[string]string
	json.Unmarshal(resp.Body.Bytes(), &loginResp)
	*token = loginResp["token"]
}
