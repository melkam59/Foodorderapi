package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"foodorderapi/internals/config"
	"foodorderapi/internals/models"
	"foodorderapi/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		db := config.DB()

		var general []*models.General

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			data := map[string]interface{}{
				"message": "Authorization header missing",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			data := map[string]interface{}{
				"message": "Invalid authorization header format",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}

		tokenString := parts[1]

		// Extract the merchant identifier from the token

		tokenparts := strings.Split(tokenString, ".")

		if len(tokenparts) != 3 {

			data := map[string]interface{}{
				"message": "Invalid token format",
			}
			return c.JSON(http.StatusUnauthorized, data)

		}

		payload, err := base64.RawURLEncoding.DecodeString(tokenparts[1])
		if err != nil {
			data := map[string]interface{}{
				"message": "Failed to decode token payload: ",
			}
			return c.JSON(http.StatusUnauthorized, data)

		}
		var claims map[string]interface{}

		err = json.Unmarshal(payload, &claims)
		if err != nil {
			data := map[string]interface{}{
				"message": "Failed to unmarshal token claims ",
			}
			return c.JSON(http.StatusUnauthorized, data)

		}

		generalID, ok := claims["sub"].(string)

		if !ok {
			data := map[string]interface{}{
				"message": "Invalid user ID",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}

		if res := db.Where("id=?", generalID).Find(&general); res.Error != nil {
			data := map[string]interface{}{
				"message": "user not found",
			}

			return c.JSON(http.StatusInternalServerError, data)
		}

		role, ok := claims["role"].(string)
		if !ok {
			data := map[string]interface{}{
				"message": "Invalid role value",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}

		publicKeyBytes := []byte(general[0].PublicKey)

		fmt.Println(general[0].PublicKey)

		publicKey, err := utils.ParseECDSAPublicKey(publicKeyBytes)
		if err != nil {
			data := map[string]interface{}{
				"message": "Invalid public key",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the ECDSA public key for verification
			return publicKey, nil
		})

		// Check for parsing errors
		if err != nil {
			fmt.Println("Failed to parse token:", err)

		}

		// Check if the token is valid
		if token.Valid {
			fmt.Println("Token is valid")

		} else {
			fmt.Println("Token is invalid")
			data := map[string]interface{}{
				"message": "Invalid token",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}
		// Token is valid, proceed with the next middleware/handler

		if role == "admin" {
			c.Set("adminID", general[0].AdminID)
			c.Set("role", role)
			fmt.Println("admin id set")
		} else {
			c.Set("merchantID", general[0].MerchantID)
			c.Set("role", role)
			fmt.Println("merchant id set")
		}
		return next(c)
	}
}

/*
func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		db := config.DB()

		var merchant []*models.Merchant

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			data := map[string]interface{}{
				"message": "Authorization header missing",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			data := map[string]interface{}{
				"message": "Invalid authorization header format",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}

		tokenString := parts[1]

		// Extract the merchant identifier from the token

		tokenparts := strings.Split(tokenString, ".")

		if len(tokenparts) != 3 {

			data := map[string]interface{}{
				"message": "Invalid token format",
			}
			return c.JSON(http.StatusUnauthorized, data)

		}

		payload, err := base64.RawURLEncoding.DecodeString(tokenparts[1])
		if err != nil {
			data := map[string]interface{}{
				"message": "Failed to decode token payload: ",
			}
			return c.JSON(http.StatusUnauthorized, data)

		}
		var claims map[string]interface{}
		err = json.Unmarshal(payload, &claims)
		if err != nil {
			data := map[string]interface{}{
				"message": "Failed to unmarshal token claims ",
			}
			return c.JSON(http.StatusUnauthorized, data)

		}

		merchantID, ok := claims["sub"].(string)
		if !ok {
			data := map[string]interface{}{
				"message": "Invalid merchant ID",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}

		if res := db.Where("id=?", merchantID).Find(&merchant); res.Error != nil {
			data := map[string]interface{}{
				"message": "merchant not found",
			}

			return c.JSON(http.StatusInternalServerError, data)
		}

		publicKeyBytes := []byte(merchant[0].PublicKey)

		fmt.Println(merchant[0].PublicKey)

		publicKey, err := utils.ParseECDSAPublicKey(publicKeyBytes)
		if err != nil {
			data := map[string]interface{}{
				"message": "Invalid public key",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}

		// Verify the token using the merchant's public key

		// Define the ECDSA algorithm
		//ecdsaAlg := jwt.SigningMethodES256

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the ECDSA public key for verification
			return publicKey, nil
		})

		// Check for parsing errors
		if err != nil {
			fmt.Println("Failed to parse token:", err)

		}

		// Check if the token is valid
		if token.Valid {
			fmt.Println("Token is valid")

		} else {
			fmt.Println("Token is invalid")
			data := map[string]interface{}{
				"message": "Invalid token",
			}
			return c.JSON(http.StatusUnauthorized, data)
		}
		// Token is valid, proceed with the next middleware/handler
		c.Set("merchantID", merchantID)
		return next(c)
	}
}


*/
