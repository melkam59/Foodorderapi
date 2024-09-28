package controllers

import (
	"fmt"
	"foodorderapi/internals/config"
	"foodorderapi/internals/models"

	"foodorderapi/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Signupadmin(c echo.Context) error {

	db := config.DB()

	var payload *models.Admin

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	hashedpassword, err := utils.Hashpassword(payload.Password)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusBadGateway, data)
	}

	private, public, err := utils.GenerateKeyPair()
	if err != nil {
		data := map[string]interface{}{
			"message": "could not generate key",
		}
		c.JSON(http.StatusBadGateway, data)
	}

	now := time.Now()

	newadmin := &models.Admin{
		AdminName:   payload.AdminName,
		Email:       payload.Email,
		Password:    hashedpassword,
		Phonenumber: payload.Phonenumber,
		PrivateKey:  private,
		PublicKey:   public,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := db.Create(&newadmin).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	newgeneraluser := &models.General{
		Name:        payload.AdminName,
		MerchantID:  "",
		AdminID:     newadmin.Id,
		Email:       payload.Email,
		Password:    hashedpassword,
		Phonenumber: payload.Phonenumber,
		Role:        "admin",
		PrivateKey:  private,
		PublicKey:   public,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := db.Create(&newgeneraluser).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	return c.JSON(http.StatusOK, newadmin)

}

func Signinadmin(c echo.Context) error {
	db := config.DB()

	var payload *models.Signininputone
	var general *models.General

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	phoneNumber := payload.Phonenumber

	if res := db.Where("phonenumber = ?", phoneNumber).Find(&general); res.Error != nil {

		data := map[string]interface{}{
			"message": "user not found",
		}

		return c.JSON(http.StatusUnauthorized, data)

	}

	password := payload.Password

	// Verify password
	passCheck := utils.VerifyPassword(general.Password, password)

	if !passCheck {
		return c.JSON(http.StatusConflict, "incorrect password or phonenumber")
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
		"user":         general,
	}

	return c.JSON(http.StatusOK, response)

}

func Logout(c echo.Context) error {
	// Retrieve the access token from the request headers or query parameters
	accessToken := c.Request().Header.Get("Authorization")
	if accessToken == "" {
		accessToken = c.QueryParam("access_token")
	}

	// Revoke the access token
	err := utils.RevokeToken(accessToken, time.Now().Add(-time.Hour*24*7))
	if err != nil {
		data := map[string]interface{}{
			"message": "Failed to log out",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	data := map[string]interface{}{
		"message": "Logged out successfully",
	}
	return c.JSON(http.StatusOK, data)
}

func UpdateAdmin(c echo.Context) error {
	db := config.DB()

	adminID := c.Param("id")

	// Retrieve the existing admin from the database
	var existingAdmin *models.Admin
	var exitinggeneraluser *models.General

	if res := db.Where("id=?", adminID).Find(&existingAdmin); res.Error != nil {
		data := map[string]interface{}{
			"message": "admin not found",
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("admin_id=?", adminID).Find(&exitinggeneraluser); res.Error != nil {
		data := map[string]interface{}{
			"message": "admin not found",
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	var payload *models.Admin

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	hashedpassword, err := utils.Hashpassword(payload.Password)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusBadGateway, data)
	}

	// Update the admin properties
	existingAdmin.AdminName = payload.AdminName
	existingAdmin.Email = payload.Email
	existingAdmin.Phonenumber = payload.Phonenumber
	existingAdmin.Password = hashedpassword
	existingAdmin.UpdatedAt = time.Now()

	// Save the updated admin to the database
	if err := db.Save(&existingAdmin).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	// Update the general table properties
	exitinggeneraluser.Name = payload.AdminName
	exitinggeneraluser.Email = payload.Email
	exitinggeneraluser.Phonenumber = payload.Phonenumber
	exitinggeneraluser.Password = hashedpassword
	exitinggeneraluser.UpdatedAt = time.Now()

	// Save the updated general table to the database
	if err := db.Save(&exitinggeneraluser).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	return c.JSON(http.StatusOK, exitinggeneraluser)
}

func DeleteAdmin(c echo.Context) error {
	db := config.DB()

	adminID := c.Param("id")

	// Retrieve the existing admin from the database
	var existingAdmin *models.Admin
	var exitinggeneraluser *models.General

	if res := db.Where("id=?", adminID).Find(&existingAdmin); res.Error != nil {
		data := map[string]interface{}{
			"message": "admin not found",
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("admin_id=?", adminID).Find(&exitinggeneraluser); res.Error != nil {
		data := map[string]interface{}{
			"message": "admin not found",
		}

		return c.JSON(http.StatusInternalServerError, data)
	}
	// Delete the admin from the database
	if err := db.Delete(&exitinggeneraluser).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}
	// Delete the admin from the database
	if err := db.Delete(&existingAdmin).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	data := map[string]interface{}{
		"message": "Admin deleted successfully",
	}
	return c.JSON(http.StatusOK, data)

}

func Signupmerchant(c echo.Context) error {

	role := c.Get("role").(string)

	// Check if the role is admin
	if role != "admin" {

		data := map[string]interface{}{

			"message": "Access denied. Only admins can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}




	db := config.DB()
	message := "could not generate key "
	var payload *models.Signupinputs

	if err := c.Bind(&payload); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	hashedpassword, err := utils.Hashpassword(payload.Password)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusBadGateway, data)
	}

	private, public, err := utils.GenerateKeyPair()
	if err != nil {
		c.JSON(http.StatusBadGateway, message)
	}

	now := time.Now()

	newmerchant := &models.Merchant{
		BusinessName:      payload.BusinessName,
		OwnerName:         payload.OwnerName,
		ContactPerson:     payload.ContactPerson,
		Phonenumber:       payload.Phonenumber,
		Email:             payload.Email,
		Password:          hashedpassword,
		MerchantShortcode: payload.MerchantShortcode,
		PrivateKey:        private,
		PublicKey:         public,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := db.Create(&newmerchant).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	newgeneraluser := &models.General{
		Name:              payload.BusinessName,
		MerchantID:        newmerchant.Id,
		AdminID:           "",
		Email:             payload.Email,
		Password:          hashedpassword,
		Phonenumber:       payload.Phonenumber,
		Role:              "merchant",
		MerchantShortcode: payload.MerchantShortcode,
		IsActive:          true,
		PrivateKey:        private,
		PublicKey:         public,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := db.Create(&newgeneraluser).Error; err != nil {

		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	merchantResponse := &models.MerchantResponse{
		Id:            newmerchant.Id,
		BusinessName:  newmerchant.BusinessName,
		OwnerName:     newmerchant.OwnerName,
		ContactPerson: newmerchant.ContactPerson,
		Email:         newmerchant.Email,
		Phonenumber:   newmerchant.Phonenumber,

		MerchantShortcode: payload.MerchantShortcode,
		PublicKey:         public,
		CreatedAt:         newmerchant.CreatedAt,
		UpdatedAt:         newmerchant.UpdatedAt,
	}

	/*
		password := strings.ReplaceAll(payload.Password, "@", "at")
		sendersms := fmt.Sprintf("You have successfully registered to Santimpay Loyalty Reward system. Your default password for first-time login is '%s'. Dont't  forget to change it later  \n Thank you for choosing us.", password)
		Statuscode, err := utils.SendSMS(sendersms, "9360", "0984006406")
		if err != nil {
			data := map[string]interface{}{
				"message": err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, data)
		}

		if Statuscode == 204 {
			return c.JSON(http.StatusAccepted, "SMS successfully sent")
		}
	*/
	return c.JSON(http.StatusOK, merchantResponse)

}

// to desplay all merchants list
func Getallmerchant(c echo.Context) error {
	role := c.Get("role").(string)

	// Check if the role is admin
	if role != "admin" {
		data := map[string]interface{}{
			"message": "Access denied. Only admins can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	db := config.DB()

	var merchants []models.Merchant

	// Retrieve all merchants
	if err := db.Find(&merchants).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	// Retrieve and attach the Is_Active attribute from the generals table
	for i := range merchants {
		var general models.General
		if err := db.Where("merchant_id = ?", merchants[i].Id).Find(&general).Error; err != nil {
			data := map[string]interface{}{
				"message": err.Error(),
			}
			return c.JSON(http.StatusInternalServerError, data)
		}
		merchants[i].IsActive = general.IsActive
		fmt.Println(general.IsActive)
	}
	fmt.Println(merchants, "those ???")
	return c.JSON(http.StatusOK, merchants)
}

func Singlemerchant(c echo.Context) error {

	role := c.Get("role").(string)

	// Check if the role is admin
	if role != "admin" {

		data := map[string]interface{}{

			"message": "Access denied. Only admins can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	db := config.DB()

	var reqBody struct {
		Phonenumber int64 `json:"phonenumber"`
	}

	var merchant []*models.Merchant

	if res := db.Where("phonenumber=?", reqBody.Phonenumber).Find(&merchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "merchant not found",
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{

		"data": merchant[0],
	}

	return c.JSON(http.StatusOK, response)

}

// to filter the merchant into Dashboardmerchat and api merchant

func UpdateMerchantbyAdmin(c echo.Context) error {
	role := c.Get("role").(string)

	// Check if the role is admin
	if role != "admin" {
		data := map[string]interface{}{
			"message": "Access denied. Only admins can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	db := config.DB()

	merchantID := c.Param("id")

	// Retrieve the existing merchant from the database
	var existingMerchant *models.Merchant
	var exitinggeneraluser *models.General

	if res := db.Where("id = ?", merchantID).Find(&existingMerchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("merchant_id=?", merchantID).Find(&exitinggeneraluser); res.Error != nil {
		data := map[string]interface{}{
			"message": "merchant not found",
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

	hashedpassword, err := utils.Hashpassword(payload.Password)

	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusBadGateway, data)
	}

	// Update the merchant properties

	existingMerchant.BusinessName = payload.BusinessName
	existingMerchant.OwnerName = payload.OwnerName
	existingMerchant.ContactPerson = payload.ContactPerson
	existingMerchant.Email = payload.Email
	existingMerchant.Phonenumber = payload.Phonenumber

	existingMerchant.MerchantShortcode = payload.MerchantShortcode
	existingMerchant.Password = hashedpassword

	existingMerchant.UpdatedAt = time.Now()

	// Save the updated merchant to the database
	if err := db.Save(&existingMerchant).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	// Update the general table properties
	exitinggeneraluser.Name = payload.BusinessName
	exitinggeneraluser.Email = payload.Email
	exitinggeneraluser.Phonenumber = payload.Phonenumber
	exitinggeneraluser.Password = hashedpassword
	exitinggeneraluser.UpdatedAt = time.Now()

	// Save the updated general table to the database
	if err := db.Save(&exitinggeneraluser).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	merchantResponse := &models.MerchantResponse{
		Id:            existingMerchant.Id,
		BusinessName:  existingMerchant.BusinessName,
		OwnerName:     existingMerchant.OwnerName,
		ContactPerson: existingMerchant.ContactPerson,
		Email:         existingMerchant.Email,
		Phonenumber:   existingMerchant.Phonenumber,

		MerchantShortcode: existingMerchant.MerchantShortcode,
		PublicKey:         existingMerchant.PublicKey,
		CreatedAt:         existingMerchant.CreatedAt,
		UpdatedAt:         existingMerchant.UpdatedAt,
	}

	return c.JSON(http.StatusOK, merchantResponse)
}

func DeleteMerchant(c echo.Context) error {
	role := c.Get("role").(string)

	// Check if the role is admin
	if role != "admin" {
		data := map[string]interface{}{
			"message": "Access denied. Only admins can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	db := config.DB()

	merchantID := c.Param("id")

	// Retrieve the existing merchant from the database
	var existingMerchant *models.Merchant
	var exitinggeneraluser *models.General

	if res := db.Where("id = ?", merchantID).Find(&existingMerchant); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if res := db.Where("merchant_id=?", merchantID).Find(&exitinggeneraluser); res.Error != nil {
		data := map[string]interface{}{
			"message": "admin not found",
		}

		return c.JSON(http.StatusInternalServerError, data)
	}
	// Delete the admin from the database
	if err := db.Delete(&exitinggeneraluser).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	// Delete the merchant from the database
	if err := db.Delete(&existingMerchant).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	data := map[string]interface{}{
		"message": "Merchant deleted successfully",
	}

	return c.JSON(http.StatusOK, data)
}

func ActivateMerchant(c echo.Context) error {

	role := c.Get("role").(string)

	// Check if the role is admin
	if role != "admin" {
		data := map[string]interface{}{
			"message": "Access denied. Only admins can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	merchantID := c.Param("id")

	// Get the merchant from the database
	db := config.DB()
	general := models.General{}
	if err := db.First(&general, "merchant_id = ?", merchantID).Error; err != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Activate the merchant
	general.IsActive = true
	if err := db.Save(&general).Error; err != nil {
		data := map[string]interface{}{
			"message": "Failed to activate merchant",
			"error":   err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	data := map[string]interface{}{
		"message":  "Merchant activated successfully",
		"merchant": general,
	}
	return c.JSON(http.StatusOK, data)
}

func DeactivateMerchant(c echo.Context) error {

	role := c.Get("role").(string)

	// Check if the role is admin
	if role != "admin" {
		data := map[string]interface{}{
			"message": "Access denied. Only admins can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	// Get the merchant ID from the request params or body
	merchantID := c.Param("id")

	// Get the merchant from the database
	db := config.DB()
	general := models.General{}
	if err := db.First(&general, "merchant_id = ?", merchantID).Error; err != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Activate the merchant
	general.IsActive = false
	if err := db.Save(&general).Error; err != nil {

		data := map[string]interface{}{
			"message": "Failed to deactivate merchant",
			"error":   err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	data := map[string]interface{}{
		"message":  "Merchant deactivated successfully",
		"merchant": general,
	}
	return c.JSON(http.StatusOK, data)
}
