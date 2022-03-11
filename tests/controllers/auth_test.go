package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ericklima-ca/bago/controllers"
	"github.com/ericklima-ca/bago/database"
	"github.com/ericklima-ca/bago/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func initDB(tb testing.TB) {
	os.Setenv("BAGO_ENV", "test")
	database.ConnectToDatabase()
}
func TestAuthLogin(t *testing.T) {
	initDB(t)
	var userTest = models.User{
		ID:       14511,
		Name:     "Erick",
		Lastname: "Lima",
		Email:    "email@email.com",
		Active:   true,
	}
	tx := database.DB.Create(&userTest)
	defer tx.Rollback()

	r := gin.Default()
	r.POST("/api/auth/login", controllers.Auth.Login)

	b, _ := json.Marshal(map[string]string{
		"login":    "14511",
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
	assert.Equal(t, true, body.Ok)
}
