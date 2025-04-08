package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	// tokenString, err := generateToken(signup.Name, signup.Id)

	// if err != nil {
	// 	models.ResponseWithError(res, http.StatusUnauthorized, "Error generating token")
	// 	return
	// }

	// response := models.AuthResponse{
	// 	Message: "Signup successfully!",
	// 	UserId:  signup.Id,
	// 	Name:    signup.Name,
	// 	Token:   tokenString,
	// }

	response, err := savingTokenAndRefresh(true, signup.Name, signup.Id)
	if err != nil {
		models.ResponseWithError(res, http.StatusUnauthorized, err.Error())
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

	response, err := savingTokenAndRefresh(false, user.Name, user.Id)
	if err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, response)

}

func AuthRefresh(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		models.ResponseWithError(res, http.StatusUnauthorized, "You don't have permission to access this resource!")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	var refreshToken models.RefreshToken
	// validate if access token exist in the database
	if err := config.Database.Preload("User").First(&refreshToken, "token = ?", tokenString).Error; err != nil {
		models.ResponseWithError(res, http.StatusUnauthorized, "You don't have permission to access this resource!")
		return
	}

	// validate refresh token if still valid else user need to login
	if _, err := validateRefreshToken(refreshToken.RefreshToken); err != nil {
		models.ResponseWithError(res, http.StatusUnauthorized, "Token expired!")
	}

	// generate new token
	tokenString, err := generateToken(refreshToken.User.Name, refreshToken.User.Id)
	if err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Error generating token!")
		return
	}

	refreshToken.Token = tokenString

	if err := config.Database.Save(&refreshToken).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to update the token: "+err.Error())
		return
	}

	response := models.AuthResponse{
		Message: "New token has been issued!",
		UserId:  refreshToken.UserId,
		Name:    refreshToken.User.Name,
		Token:   tokenString,
	}

	models.ResponseWithJSON(res, http.StatusCreated, response)
}

func AuthCurrent(res http.ResponseWriter, req *http.Request) {

}

// private functions
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

func generateRefreshToken(expirationTime time.Time, name string, userId uuid.UUID) (string, error) {
	claims := &models.JwtClaims{
		Name: name,
		Id:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.RefreshJWTKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func savingTokenAndRefresh(isSignup bool, name string, userId uuid.UUID) (models.AuthResponse, error) {
	var response models.AuthResponse

	tokenString, err := generateToken(name, userId)
	if err != nil {
		return response, fmt.Errorf("error generating token: %w", err)
	}

	var refreshToken models.RefreshToken

	result := config.Database.Where("user_id = ?", userId).Delete(&refreshToken)
	if result.Error != nil {
		return response, fmt.Errorf("error deleting user tokens: %w", result.Error)
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days expiration
	refreshTokenString, err := generateRefreshToken(expirationTime, name, userId)
	if err != nil {
		return response, fmt.Errorf("error generating refresh token: %w", err)
	}

	refreshToken = models.RefreshToken{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
		UserId:       userId,
		ExpiresAt:    expirationTime,
	}

	result = config.Database.Create(&refreshToken)
	if result.Error != nil {
		return response, fmt.Errorf("error saving refresh token: %w", result.Error)
	}

	message := "Login successfully!"
	if isSignup {
		message = "Signup successful!"
	}

	response = models.AuthResponse{
		Message: message,
		UserId:  userId,
		Name:    name,
		Token:   tokenString,
	}

	return response, nil
}

func validateRefreshToken(tokenString string) (string, error) {
	claims := &models.JwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.RefreshJWTKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return tokenString, fmt.Errorf("token expired")
		}
		return tokenString, err
	}

	if !token.Valid {
		return tokenString, fmt.Errorf("token expired")
	}
	return tokenString, nil
}
