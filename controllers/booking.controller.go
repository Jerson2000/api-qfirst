package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
)

func BookingCreate(res http.ResponseWriter, req *http.Request) {
	var booking models.Booking
	err := json.NewDecoder(req.Body).Decode(&booking)
	if err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Invalid request payload!\t"+err.Error())
		return
	}

	claims, ok := req.Context().Value("claims").(*models.JwtClaims)
	if !ok || claims == nil {
		models.ResponseWithError(res, http.StatusUnauthorized, "Invalid token or claims")
		return
	}

	userId := claims.Id
	if userId != booking.UserId {
		models.ResponseWithError(res, http.StatusForbidden, "You don't have permission!")
		return
	}

	// Check if the user exists
	var user models.User
	if err := config.Database.First(&user, "id=?", userId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "User not found!\t"+err.Error())
		return
	}

	// Check if the service exists
	var service models.Services
	if err := config.Database.First(&service, booking.ServiceId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Service not found!\t"+err.Error())
		return
	}

	// Check if there's already an existing booking for the same user and service with overlapping dates
	var existingBooking models.Booking
	err = config.Database.Where("user_id = ? AND service_id = ? AND ((start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?))",
		userId, booking.ServiceId, booking.StartDate, booking.EndDate, booking.StartDate, booking.EndDate).
		First(&existingBooking).Error

	if err == nil {
		// If a booking is found, it means there's a conflict (double booking)
		models.ResponseWithError(res, http.StatusConflict, "You already have a booking for this service during this time!")
		return
	}

	// If no conflicting booking, create the new booking
	result := config.Database.Create(&booking)
	if result.Error != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, result.Error.Error())
		return
	}

	models.ResponseWithJSON(res, http.StatusCreated, map[string]string{"message": "Booking has been created!"})
}

func BookingUpdate(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	bookingId := vars["id"]

	var bookingBody models.Booking
	err := json.NewDecoder(req.Body).Decode(&bookingBody)

	if err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Invalid request payload!")
		return
	}

	// check is exist
	var booking models.Services
	if err := config.Database.First(&booking, bookingId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Booking not found!\t"+err.Error())
		return
	}

	// check is exist
	var user models.User
	if err := config.Database.First(&user, "id=?", bookingBody.UserId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "User not found!\t"+err.Error())
		return
	}

	// check is exist
	var service models.Services
	if err := config.Database.First(&service, bookingBody.ServiceId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Service not found!\t"+err.Error())
		return
	}

	// updating
	if err := config.Database.Model(&service).Updates(bookingBody).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to update booking!")
		return
	}
	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"message": "Booking has been updated!"})
}

func BookingList(res http.ResponseWriter, req *http.Request) {

	var booking []models.Booking
	if err := config.Database.Find(&booking).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": booking})
}

func BookingGetById(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	bookingId := vars["id"]

	var booking models.Booking
	if err := config.Database.First(&booking, bookingId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Booking not found!")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": booking})
}

func BookingDelete(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	bookingId := vars["id"]

	var booking models.Booking
	if err := config.Database.First(&booking, bookingId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "booking not found!")
		return
	}

	if err := config.Database.Delete(&booking).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to delete booking!")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"message": "Booking has been deleted successfully!"})
}
