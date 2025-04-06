package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
	"golang.org/x/crypto/bcrypt"
)

func UserCreate(res http.ResponseWriter, req *http.Request) {
	var user models.User

	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Invalid request payload!")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Error hashing the password")
		return
	}

	// set the password
	user.Password = string(hashedPassword)

	result := config.Database.Create(&user)

	if result.Error != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Invalid request payload!")
		return
	}

	models.ResponseWithJSON(res, http.StatusCreated, map[string]string{"message": "User has been created!"})
}

func UserUpdate(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["id"]

	var userBody models.User

	err := json.NewDecoder(req.Body).Decode(&userBody)

	if err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Invalid request payload!")
		return
	}

	var user models.User
	if err := config.Database.First(&user, userId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "User not found!")
		return
	}

	if err := config.Database.Model(&user).Updates(userBody).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to update user!")
		return
	}
	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"message": "User has been updated!"})
}

func UserList(res http.ResponseWriter, req *http.Request) {

	claims := req.Context().Value("claims").(*models.JwtClaims)
	userId := claims.ID

	if len(userId) == 0 {
		models.ResponseWithError(res, http.StatusUnauthorized, "You don't permission!")
		return
	}

	var users []models.User
	if err := config.Database.Find(&users).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": users})
}

func UserGetById(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["id"]

	var user models.User
	if err := config.Database.First(&user, userId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "User not found!")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": user})
}

func UserDelete(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["id"]

	var user models.User
	if err := config.Database.First(&user, userId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "User not found!")
		return
	}

	if err := config.Database.Delete(&user).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to update user!")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"message": "User has been deleted successfully!"})
}
