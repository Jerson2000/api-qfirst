package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
)

// BookingCreate godoc
// @Summary Create a new booking.
// @Description Create a new booking for a user. The booking will be stored in the database, and the user will receive a confirmation.
// @Tags Booking
// @Accept  json
// @Produce  json
// @Param booking body models.Booking true "Booking details"
// @Success 201 {object} models.Booking "Booking created successfully"
// @Failure 400 {object} map[string]string "Invalid input or missing data"
// @Failure 404 {object} map[string]string "Service or User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /bookings [post]
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

// BookingUpdate godoc
// @Summary Update an existing booking.
// @Description Update an existing booking by its ID. The booking details will be updated, and the changes will be saved to the database.
// @Tags Booking
// @Accept  json
// @Produce  json
// @Param id path uint true "Booking ID"
// @Param booking body models.Booking true "Booking details to update"
// @Success 200 {object} models.Booking "Booking updated successfully"
// @Failure 400 {object} map[string]string "Invalid input or missing data"
// @Failure 404 {object} map[string]string "Booking not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /bookings/{id} [put]
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

// BookingList godoc
// @Summary Get a list of bookings.
// @Description Retrieve a list of bookings, with optional filtering by user ID, service ID, or status.
// @Tags Booking
// @Accept  json
// @Produce  json
// @Param user_id query string false "User ID"
// @Param service_id query string false "Service ID"
// @Param status query string false "Booking Status" Enums(pending, confirmed, cancelled)
// @Param page query int false "Page number for pagination"
// @Param page_size query int false "Page size for pagination"
// @Success 200 {array} models.Booking "List of bookings"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /bookings [get]
func BookingList(res http.ResponseWriter, req *http.Request) {
	var booking []models.Booking

	// Get query parameters
	userID := req.URL.Query().Get("user_id")
	serviceID := req.URL.Query().Get("service_id")
	status := req.URL.Query().Get("status")
	page := req.URL.Query().Get("page")
	pageSize := req.URL.Query().Get("page_size")

	// Build query conditions dynamically based on query parameters
	query := config.Database.Model(&models.Booking{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Handle pagination if provided
	if page != "" && pageSize != "" {
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			models.ResponseWithError(res, http.StatusBadRequest, "Invalid page number")
			return
		}
		pageSizeNum, err := strconv.Atoi(pageSize)
		if err != nil {
			models.ResponseWithError(res, http.StatusBadRequest, "Invalid page size")
			return
		}

		// Implement pagination (Skip and limit)
		offset := (pageNum - 1) * pageSizeNum
		query = query.Offset(offset).Limit(pageSizeNum)
	}

	// Fetch the bookings with the built query
	if err := query.Find(&booking).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the response with bookings
	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": booking})
}

// BookingGetById godoc
// @Summary Get a booking by ID.
// @Description Retrieve a booking by its unique ID.
// @Tags Booking
// @Accept  json
// @Produce  json
// @Param id path uint true "Booking ID"
// @Success 200 {object} models.Booking "Booking found"
// @Failure 400 {object} map[string]string "Invalid Booking ID"
// @Failure 404 {object} map[string]string "Booking not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /bookings/{id} [get]
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

// BookingDelete godoc
// @Summary Delete a booking by ID.
// @Description Delete a booking from the platform using its unique ID.
// @Tags Booking
// @Accept  json
// @Produce  json
// @Param id path uint true "Booking ID"
// @Success 204 "Booking deleted successfully"
// @Failure 400 {object} map[string]string "Invalid Booking ID"
// @Failure 404 {object} map[string]string "Booking not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /bookings/{id} [delete]
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
