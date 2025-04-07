package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/models"
)

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

func ServiceList(res http.ResponseWriter, req *http.Request) {

	var services []models.Services
	if err := config.Database.Find(&services).Error; err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, err.Error())
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]any{"data": services})
}

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
