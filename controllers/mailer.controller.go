package controllers

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/jerson2000/api-qfirst/config"
	"github.com/jerson2000/api-qfirst/mailer"
	"github.com/jerson2000/api-qfirst/models"
)

type requestBody struct {
	Email string `json:"email"`
	Code  string `json:"code,omitempty"`
}

func MailerGenerateOTP(res http.ResponseWriter, req *http.Request) {
	var reqBody requestBody

	err := json.NewDecoder(req.Body).Decode(&reqBody)

	if err != nil {
		models.ResponseWithError(res, http.StatusBadRequest, "invalid payload")
		return
	}

	// find user by email
	var user models.User
	if err := config.Database.First(&user, "email = ?", reqBody.Email).Error; err != nil {
		models.ResponseWithError(res, http.StatusNotFound, "user not found")
		return
	}

	otpCode, err := generateCode()
	if err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "error generating OTP Code")
		return
	}
	var otp models.OTP

	// delete all user otp
	result := config.Database.Where("user_id = ?", user.Id).Delete(&otp)
	if result.Error != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "error deleting user otp")
		return
	}

	otp = models.OTP{
		Code:      otpCode,
		Timestamp: time.Now(),
		Expiry:    time.Now().Add(5 * time.Minute), // 5mins. expiry
		UserId:    user.Id,
	}

	// send to email
	toMail := user.Email
	err = mailer.SendOTPMail(toMail, strconv.Itoa(otpCode))
	if err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "error sending OTP to email")
		return
	}

	// save otp
	result = config.Database.Create(&otp)
	if result.Error != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "error saving user otp")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]string{"message": "OTP code has been send to your email valid for 5 minutes, kindly check it."})

}

func MailerValidateOTP(res http.ResponseWriter, req *http.Request) {

	var reqBody requestBody

	err := json.NewDecoder(req.Body).Decode(&reqBody)

	if err != nil || reqBody.Code == "" {
		models.ResponseWithError(res, http.StatusBadRequest, "invalid payload")
		return
	}

	// find user by email
	var user models.User
	if err := config.Database.First(&user, "email = ?", reqBody.Email).Error; err != nil {
		models.ResponseWithError(res, http.StatusNotFound, "user not found")
		return
	}

	var otp models.OTP

	if err := config.Database.First(&otp, "user_id = ?", user.Id).Error; err != nil {
		models.ResponseWithError(res, http.StatusNotFound, "user otp not found")
		return
	}

	reqCode, err := strconv.Atoi(reqBody.Code)
	if err != nil {
		models.ResponseWithError(res, http.StatusInternalServerError, "error parsing the otp to int")
		return
	}

	if otp.Code != reqCode {
		models.ResponseWithError(res, http.StatusNotFound, "OTP code is not valid")
		return
	}

	if time.Now().After(otp.Expiry) {
		models.ResponseWithError(res, http.StatusNotFound, "OTP code is not valid")
		return
	}

	models.ResponseWithJSON(res, http.StatusOK, map[string]string{"message": "validated"})
}

func generateCode() (int, error) {
	// Generate a random number between 100000 and 999999
	max := big.NewInt(900000)
	min := big.NewInt(100000)
	num, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return int(num.Int64()) + int(min.Int64()), nil
}
