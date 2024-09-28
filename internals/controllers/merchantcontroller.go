package controllers

import (
	"errors"
	"foodorderapi/internals/config"
	"foodorderapi/internals/models"

	"foodorderapi/utils"
	"net/http"

	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Signin(c echo.Context) error {

	db := config.DB()
	var payload *models.Merchantsignin
	var general *models.General
	var merchant *models.Merchant

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("merchant_shortcode = ?", payload.MerchantShortcode).Find(&merchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "you are not allowed to use this system! Please Register first",
		}

		return c.JSON(http.StatusForbidden, data)
	}

	if res := db.Where("merchant_shortcode = ?", payload.MerchantShortcode).Find(&general); res.Error != nil {

		data := map[string]interface{}{
			"message": "you are not allowed to use this system! Please Register first",
		}

		return c.JSON(http.StatusForbidden, data)

	}

	if !general.IsActive {
		return c.JSON(http.StatusForbidden, "you are currently inactive Please contact the support team")
	}

	password := payload.Password

	// Verify password
	passCheck := utils.VerifyPassword(general.Password, password)
	passCheck2 := utils.VerifyPassword(merchant.Password, password)

	if !passCheck || !passCheck2 {
		return c.JSON(http.StatusConflict, "incorrect password! Please try again")
	}

	ttl := 24 * time.Hour

	access_token, err := utils.Createtoken(ttl, general.Id, general.Role, []byte(general.PrivateKey))
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, data)
	}

	response := map[string]interface{}{
		"access_token": access_token,
		"user":         merchant,
	}

	return c.JSON(http.StatusOK, response)
}

func Me(c echo.Context) error {

	db := config.DB()
	merchantID := c.Get("merchantID").(string)
	var merchant []models.Merchant

	if res := db.Where("id=?", merchantID).Find(&merchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)

	}

	response := map[string]interface{}{

		"user": merchant,
	}

	return c.JSON(http.StatusOK, response)

}

func UpdateMerchant(c echo.Context) error {
	db := config.DB()
	merchantID := c.Get("merchantID").(string)

	var existingMerchant *models.Merchant
	var existingGeneralUser *models.General

	if res := db.Where("id = ?", merchantID).Find(&existingMerchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("merchant_id=?", merchantID).Find(&existingGeneralUser); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	var payload *models.Merchant

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	hashedPassword, err := utils.Hashpassword(payload.Password)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadGateway, data)
	}

	if payload.BusinessName != "" {
		existingMerchant.BusinessName = payload.BusinessName
		existingGeneralUser.Name = payload.BusinessName
	}
	if payload.OwnerName != "" {
		existingMerchant.OwnerName = payload.OwnerName
	}
	if payload.ContactPerson != "" {
		existingMerchant.ContactPerson = payload.ContactPerson
	}
	if payload.Email != "" {
		existingMerchant.Email = payload.Email
		existingGeneralUser.Email = payload.Email
	}
	if payload.Phonenumber != 0 {
		existingMerchant.Phonenumber = payload.Phonenumber
		existingGeneralUser.Phonenumber = payload.Phonenumber
	}
	if payload.Password != "" {
		existingMerchant.Password = hashedPassword
		existingGeneralUser.Password = hashedPassword
	}

	if payload.MerchantShortcode != 0 {

		existingMerchant.MerchantShortcode = existingMerchant.MerchantShortcode
		existingGeneralUser.MerchantShortcode = existingMerchant.MerchantShortcode

	}

	existingMerchant.IsUpdated = true
	existingMerchant.UpdatedAt = time.Now()
	existingGeneralUser.UpdatedAt = time.Now()

	if err := db.Save(&existingMerchant).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if err := db.Save(&existingGeneralUser).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	merchantResponse := &models.MerchantResponse{
		Id:                existingMerchant.Id,
		BusinessName:      existingMerchant.BusinessName,
		OwnerName:         existingMerchant.OwnerName,
		ContactPerson:     existingMerchant.ContactPerson,
		Email:             existingMerchant.Email,
		Phonenumber:       existingMerchant.Phonenumber,
		MerchantShortcode: existingMerchant.MerchantShortcode,
		PublicKey:         existingMerchant.PublicKey,
		IsUpdated:         existingMerchant.IsUpdated,
		CreatedAt:         existingMerchant.CreatedAt,
		UpdatedAt:         existingMerchant.UpdatedAt,
	}

	return c.JSON(http.StatusOK, merchantResponse)
}

// for forget password

func Forgetpassword(c echo.Context) error {
	db := config.DB()

	var payload *models.Signininputs
	var merchant *models.Merchant

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	MerchantShortcode := payload.MerchantShortcode

	if MerchantShortcode == 0 {
		data := map[string]interface{}{
			"message": "Please fill out your phone number.",
		}
		return c.JSON(http.StatusUnauthorized, data)
	}

	if err := db.Where("merchant_shortcode = ?", MerchantShortcode).First(&merchant).Error; err != nil {
		// Handle query execution error
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return c.JSON(http.StatusUnauthorized, "You are not allowed to use this system! Please register first.")
		}

		// Handle other query errors
		data := map[string]interface{}{
			"message": "An error occurred while querying the database.",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{

		"user": merchant,
	}

	return c.JSON(http.StatusOK, response)

}

func GetMerchantByShortCode(c echo.Context) error {
	db := config.DB()
	var merchant models.Merchant

	var reqBody struct {
		MerchantShortcode int64 `json:"merchantshortcode"`
	}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if res := db.Where("merchant_shortcode=?", reqBody.MerchantShortcode).Find(&merchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)

	}

	return c.JSON(http.StatusOK, merchant)
}
