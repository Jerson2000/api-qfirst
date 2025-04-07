package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthLoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AuthSignup(res http.ResponseWriter, req *http.Request) {
	var signup models.User

	err := json.NewDecoder(req.Body).Decode(&signup)

	if err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Invalid request payload!")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)

	if err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Error hashing the password")
		return
	}

	// set the password
	signup.Password = string(hashedPassword)

	result := config.Database.Create(&signup)

	if result.Error != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to signup.")
		return
	}

	tokenString, err := generateToken(signup.Name, signup.Id)

	if err != nil {
		models.ResponseWithError(res, http.StatusUnauthorized, "Error generating token")
		return
	}

	response := models.AuthResponse{
		Message: "Signup successfully!",
		UserId:  signup.Id,
		Name:    signup.Name,
		Token:   tokenString,
	}

	models.ResponseWithJSON(res, http.StatusCreated, response)
}

func AuthLogin(res http.ResponseWriter, req *http.Request) {
	var login AuthLoginStruct
	err := json.NewDecoder(req.Body).Decode(&login)

	if err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Invalid request payload!")
		return
	}

	var user models.User

	if err := config.Database.First(&user, "email = ?", login.Email).Error; err != nil {
		models.ResponseWithError(res, http.StatusUnauthorized, "Incorrect credentials!")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))

	if err != nil {
		models.ResponseWithError(res, http.StatusUnauthorized, "Incorrect credentials!")
		return
	}

	tokenString, err := generateToken(user.Name, user.Id)

	if err != nil {
		models.ResponseWithError(res, http.StatusUnauthorized, "Error generating token")
		return
	}

	response := models.AuthResponse{
		Message: "Login successfully!",
		UserId:  user.Id,
		Name:    user.Name,
		Token:   tokenString,
	}

	models.ResponseWithJSON(res, http.StatusCreated, response)

}

func AuthRefresh(res http.ResponseWriter, req *http.Request) {

}

func AuthCurrent(res http.ResponseWriter, req *http.Request) {

}

func generateToken(name string, userId uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &models.JwtClaims{
		Name: name,
		Id:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
