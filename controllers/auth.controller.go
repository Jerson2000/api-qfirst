package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/csrf"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthLoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup Authentication godoc
// @Summary Signing up to the platform.
// @Description Add a new user to the platform. The request must include valid user data, including name, email, and password. The user will be created and saved in the database.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param user body models.User true "User data to sign up"
// @Success 201 {object} models.User "User created successfully"
// @Failure 400 {object} map[string]string "Invalid input or missing data"
// @Failure 409 {object} map[string]string "Email already in use"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/signup [post]
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

	response, err := savingTokenAndRefresh(true, signup)
	if err != nil {
		models.ResponseWithError(res, http.StatusUnauthorized, err.Error())
	}

	models.ResponseWithJSON(res, http.StatusCreated, response)
}

// AuthLogin godoc
// @Summary Log in to the platform.
// @Description Authenticate a user by verifying their email and password. On successful authentication, a JWT token will be generated for the user.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param login body AuthLoginStruct true "User login credentials"
// @Success 200 {object} map[string]string "Login successful"
// @Failure 400 {object} map[string]string "Invalid input or missing data"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/login [post]
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

	response, err := savingTokenAndRefresh(false, user)
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
	tokenString, err := generateToken(refreshToken.User)
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
	claims := req.Context().Value("claims").(*models.JwtClaims)
	models.ResponseWithJSON(res, http.StatusOK, claims)
}
func AuthRequestCSRFToken(res http.ResponseWriter, req *http.Request) {
	models.ResponseWithJSON(res, http.StatusOK, map[string]string{"message": csrf.Token(req)})
}

// private functions
func generateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &models.JwtClaims{
		Name: user.Name,
		Id:   user.Id,
		Role: *user.Role,
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

func generateRefreshToken(expirationTime time.Time, user models.User) (string, error) {
	claims := &models.JwtClaims{
		Name: user.Name,
		Id:   user.Id,
		Role: *user.Role,
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

func savingTokenAndRefresh(isSignup bool, user models.User) (models.AuthResponse, error) {
	var response models.AuthResponse

	tokenString, err := generateToken(user)
	if err != nil {
		return response, fmt.Errorf("error generating token: %w", err)
	}

	var refreshToken models.RefreshToken

	result := config.Database.Where("user_id = ?", user.Id).Delete(&refreshToken)
	if result.Error != nil {
		return response, fmt.Errorf("error deleting user tokens: %w", result.Error)
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days expiration
	refreshTokenString, err := generateRefreshToken(expirationTime, user)
	if err != nil {
		return response, fmt.Errorf("error generating refresh token: %w", err)
	}

	refreshToken = models.RefreshToken{
		Token:        tokenString,
		RefreshToken: refreshTokenString,
		UserId:       user.Id,
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
		UserId:  user.Id,
		Name:    user.Name,
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
