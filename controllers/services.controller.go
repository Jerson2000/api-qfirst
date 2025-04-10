package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
)

// ServiceCreate godoc
// @Summary Create a new service.
// @Description Add a new service to the platform.
// @Tags Service
// @Accept  json
// @Produce  json
// @Param service body models.Services true "Service data"
// @Success 201 {object} models.Services "Service created successfully"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /services [post]
func ServiceCreate(res http.ResponseWriter, req *http.Request) {
	var service models.Services

	err := json.NewDecoder(req.Body).Decode(&service)

	if err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Invalid request payload!")
		return
	}

	result := config.Database.Create(&service)
	if result.Error != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, result.Error.Error())
		return
	}

	models.ResponseWithJSON(res, http.StatusCreated, map[string]string{"message": "Service has been created!"})
}

// ServiceUpdate godoc
// @Summary Update an existing service.
// @Description Update a service's details using its unique ID.
// @Tags Service
// @Accept  json
// @Produce  json
// @Param id path uint true "Service ID"
// @Param service body models.Services true "Updated service data"
// @Success 200 {object} models.Services "Service updated successfully"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 404 {object} map[string]string "Service not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /services/{id} [put]
func ServiceUpdate(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	serviceId := vars["id"]

	var serviceBody models.Services

	err := json.NewDecoder(req.Body).Decode(&serviceBody)

	if err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Invalid request payload!\t"+err.Error())
		return
	}

	var service models.Services
	if err := config.Database.First(&service, serviceId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Service not found!\t"+err.Error())
		return
	}

	if err := config.Database.Model(&service).Updates(serviceBody).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to update service!")
		return
	}
	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"message": "Service has been updated!"})
}

// ServiceList godoc
// @Summary Get all services.
// @Description Retrieve a list of all available services in the system.
// @Tags Service
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Services "List of services"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /services [get]
func ServiceList(res http.ResponseWriter, req *http.Request) {

	var services []models.Services
	if err := config.Database.Find(&services).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": services})
}

// ServiceGetById godoc
// @Summary Get service by ID.
// @Description Retrieve a service's details by its unique ID.
// @Tags Service
// @Accept  json
// @Produce  json
// @Param id path uint true "Service ID"
// @Success 200 {object} models.Services "Service details"
// @Failure 400 {object} map[string]string "Invalid Service ID"
// @Failure 404 {object} map[string]string "Service not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /services/{id} [get]
func ServiceGetById(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	serviceId := vars["id"]

	var service models.Services
	if err := config.Database.First(&service, serviceId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Service not found!")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": service})
}

// ServiceDelete godoc
// @Summary Delete service by ID.
// @Description Delete a service from the system by its unique ID.
// @Tags Service
// @Accept  json
// @Produce  json
// @Param id path uint true "Service ID"
// @Success 200 {object} map[string]string "Service deleted successfully"
// @Failure 400 {object} map[string]string "Invalid Service ID"
// @Failure 404 {object} map[string]string "Service not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /services/{id} [delete]
func ServiceDelete(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	serviceId := vars["id"]

	var service models.Services
	if err := config.Database.First(&service, serviceId).Error; err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "Service not found!")
		return
	}

	if err := config.Database.Delete(&service).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "Failed to delete service!")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"message": "Service has been deleted successfully!"})
}
