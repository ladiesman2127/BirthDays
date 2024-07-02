package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/ladiesman2127/birthdays/internal/app"
	"github.com/ladiesman2127/birthdays/internal/app/models"
)

var userSession models.Session

func initTest() (*gin.Engine, *httptest.ResponseRecorder) {
	godotenv.Load("../.env")
	app := app.New()
	w := httptest.NewRecorder()
	return app, w
}

func TestSignUp(t *testing.T) {
	app, w := initTest()
	user := "user"
	password := "password"
	phone := "+7 999 999 99 99"
	name := "some name"
	birthday := "01.07.2023"
	newUser := models.User{
		Login:    &user,
		Password: &password,
		Phone:    &phone,
		Name:     &name,
		BirthDay: &birthday,
	}

	newUserJson, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/api/signup", strings.NewReader(string(newUserJson)))
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSignUpDuplicate(t *testing.T) {
	app, w := initTest()
	login := "user"
	password := "password"
	phone := "+7 999 999 99 99"
	name := "some name"
	birthday := "01.07.2023"
	newUser := models.NewUser(&login, &password, &name, &birthday, &phone)

	newUserJson, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/api/signup", strings.NewReader(string(newUserJson)))
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestAuthCorrectCredentials(t *testing.T) {
	app, w := initTest()
	login := "user"
	password := "password"
	credentials := models.NewCredentials(&login, &password)

	credentialsJson, _ := json.Marshal(credentials)
	req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(credentialsJson)))
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	session := models.Session{}
	json.Unmarshal(w.Body.Bytes(), &session)
}

func TestAuthIncorrectCredentials(t *testing.T) {
	app, w := initTest()
	login := "user"
	password := "passwordasd"
	credentials := models.NewCredentials(&login, &password)

	credentialsJson, _ := json.Marshal(credentials)
	req, _ := http.NewRequest("POST", "/api/auth", strings.NewReader(string(credentialsJson)))
	app.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
