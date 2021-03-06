package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ericklima-ca/bago/controllers"
	"github.com/ericklima-ca/bago/database"
	"github.com/ericklima-ca/bago/models"
	"github.com/ericklima-ca/bago/router"
	cachingservice "github.com/ericklima-ca/bago/services/caching_service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setUp(t *testing.T) (*testing.T, *gorm.DB, controllers.AuthController) {
	t.Setenv("BAGO_ENV", "test")
	dbs := database.DatabaseServer{
		Models: []interface{}{&models.User{}},
	}
	db, err := dbs.Connect()
	if err != nil {
		log.Fatalf("DB not connected: %v", err.Error())
	}
	authController := controllers.AuthController{
		DB: db,
	}

	return t, db, authController
}
func TestAuthLoginSucessAndTokenCreated(t *testing.T) {
	_, db, auth := setUp(t)

	var userTest = models.User{
		ID:       14511,
		Name:     "Erick",
		Lastname: "Lima",
		Email:    "email@email.com",
		Active:   true,
	}
	var userFormTest = models.UserFormData{
		User:     userTest,
		Password: "123456",
	}
	tx := db.Create(userFormTest.GetUser)
	t.Cleanup(func() {
		tx.Rollback()
	})

	r := gin.Default()
	r.POST("/api/auth/login", auth.Login)

	b, _ := json.Marshal(map[string]string{
		"id":        "14511",
		"password1": "123456",
	})

	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var body struct {
		Ok   bool
		Body map[string]interface{}
	}

	json.Unmarshal(res.Body.Bytes(), &body)
	assert.Equal(t, false, body.Ok)
}

func TestAuthLoginFail(t *testing.T) {
	_, db, auth := setUp(t)

	var userTest = models.User{
		ID:       14511,
		Name:     "Erick",
		Lastname: "Lima",
		Email:    "email@email.com",
		Active:   true,
	}
	var userFormTest = models.UserFormData{
		User:     userTest,
		Password: "123456",
	}
	tx := db.Create(userFormTest.GetUser)
	t.Cleanup(func() {
		tx.Rollback()
	})

	r := gin.Default()
	r.POST("/api/auth/login", auth.Login)

	b, _ := json.Marshal(map[string]string{
		"id":       "14511",
		"password": "123456",
	})

	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var body struct {
		Ok   bool
		Body map[string]interface{}
	}

	json.Unmarshal(res.Body.Bytes(), &body)
	log.Println(body)
	assert.NotEqual(t, true, body.Ok)
}

func TestAuthLoginUserNotActive(t *testing.T) {
	_, db, auth := setUp(t)

	var userTest = models.User{
		ID:       14511,
		Name:     "Erick",
		Lastname: "Lima",
		Email:    "email@email.com",
	}
	var userFormTest = models.UserFormData{
		User:     userTest,
		Password: "123456",
	}
	tx := db.Create(userFormTest.GetUser())
	t.Cleanup(func() {
		tx.Rollback()
	})

	r := gin.Default()
	r.POST("/api/auth/login", auth.Login)

	b, _ := json.Marshal(map[string]interface{}{
		"id":       14511,
		"password": "123456",
	})

	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var body struct {
		Ok   bool
		Body map[string]interface{}
	}

	json.Unmarshal(res.Body.Bytes(), &body)
	assert.Equal(t, "user not active", body.Body["error"])
}

func TestAuthLoginFailPayload(t *testing.T) {
	_, _, auth := setUp(t)

	r := gin.Default()
	r.POST("/api/auth/login", auth.Login)

	b, _ := json.Marshal(map[string]string{
		"wrong_login": "14511",
		"password":    "123456",
	})

	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var body struct {
		Ok   bool
		Body map[string]interface{}
	}

	json.Unmarshal(res.Body.Bytes(), &body)
	assert.Equal(t, "invalid payload", body.Body["error"])
}

func TestSignupSucess(t *testing.T) {
	_, db, _ := setUp(t)

	authController := controllers.AuthController{
		DB: db,
	}
	routerServer := router.Router{
		AuthController: &authController,
	}
	server := routerServer.LoadRoutes()

	payloadSignup := gin.H{
		"id":       14512,
		"name":     "Erick",
		"lastname": "Amorim",
		"password": "123456789",
		"email":    "email@email.com",
	}
	jsonBody, _ := json.Marshal(payloadSignup)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(res, req)

	var body struct {
		Ok   bool
		Body map[string]interface{}
	}

	str := cachingservice.GetToken("signup", 14512)
	json.Unmarshal(res.Body.Bytes(), &body)
	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Equal(t, "user created", body.Body["msg"])
	assert.NotEqual(t, str, "", "empty token")
}

func TestSignupPayloadFail(t *testing.T) {
	_, db, _ := setUp(t)

	authController := controllers.AuthController{
		DB: db,
	}
	routerServer := router.Router{
		AuthController: &authController,
	}
	server := routerServer.LoadRoutes()

	payloadSignup := gin.H{
		"id":       14511,
		"name":     "Erick",
		"lastname": "Amorim",
		"password": 123456789, // as int
		"email":    "email@email.com",
	}
	jsonBody, _ := json.Marshal(payloadSignup)

	res := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/auth/signup", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(res, req)

	var body struct {
		Ok   bool
		Body map[string]interface{}
	}

	json.Unmarshal(res.Body.Bytes(), &body)

	assert.Equal(t, http.StatusBadRequest, res.Code)
	assert.Equal(t, "data does not match", body.Body["error"])
}
