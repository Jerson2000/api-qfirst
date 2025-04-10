package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserCreate godoc
// @Summary Create a new user
// @Description Add a new user to the platform
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.User true "User data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /users [post]
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
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			models.ResponseWithError(res, http.StatusBadRequest, "Email already exist!")
			return
		}
		models.ResponseWithError(res, http.StatusBadRequest, result.Error.Error())
		return
	}

	models.ResponseWithJSON(res, http.StatusCreated, map[string]string{"message": "User has been created!"})
}

// UserUpdate godoc
// @Summary Update a user
// @Description Update user information by ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /users/{id} [put]
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
	if err := config.Database.First(&user, "id=?", userId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "User not found!")
		return
	}

	if err := config.Database.Model(&user).Updates(userBody).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to update user!")
		return
	}
	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"message": "User has been updated!"})
}

// UserList godoc
// @Summary List all users
// @Description Get a list of all users
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /users [get]
func UserList(res http.ResponseWriter, req *http.Request) {

	// claims := req.Context().Value("claims").(*models.JwtClaims)
	// userId := claims.Id

	// if len(userId) == 0 {
	// 	models.ResponseWithError(res, http.StatusUnauthorized, "You don't have permission!")
	// 	return
	// }

	var users []models.User
	if err := config.Database.Find(&users).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": users})
}

// UserGetById godoc
// @Summary Get user by ID
// @Description Retrieve a user's details by their unique ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /users/{id} [get]
func UserGetById(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["id"]

	var user models.User
	if err := config.Database.First(&user, "id=?", userId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "User not found!")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": user})
}

// UserDelete godoc
// @Summary Delete user by ID
// @Description Permanently deletes a user by their unique ID
// @Tags User
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /users/{id} [delete]
func UserDelete(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["id"]

	var user models.User
	if err := config.Database.First(&user, "id=?", userId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "User not found!")
		return
	}

	if err := config.Database.Delete(&user).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to delete service!")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"message": "User has been deleted successfully!"})
}
