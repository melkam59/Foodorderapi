package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"foodorderapi/internals/config"
	"foodorderapi/internals/models"

	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateKeyPair() (string, string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

func ParseECDSAPrivateKeyFromPEM(pemData []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("failed to decode PEM data")
	}

	var privKey *ecdsa.PrivateKey
	var err error

	if block.Type == "EC PRIVATE KEY" {
		privKey, err = x509.ParseECPrivateKey(block.Bytes)
	} else {
		privKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		ecdsaPrivKey, ok := privKeyInterface.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.New("parsed key is not an ECDSA private key")
		}

		privKey = ecdsaPrivKey
	}

	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func ParseECDSAPublicKey(publicKeyPEM []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, errors.New("failed to decode PEM data")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("parsed key is not an ECDSA public key")
	}

	return ecdsaPubKey, nil

}

func Createtoken(ttl time.Duration, logged_in interface{}, role interface{}, privateKeyBytes []byte) (string, error) {

	privateKey, err := ParseECDSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	expirationTime := now.Add(ttl).Unix()

	claims := jwt.MapClaims{
		"sub":  logged_in,
		"exp":  expirationTime,
		"role": role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// Generate ECDSA signature
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return tokenString, nil
}

func RevokeToken(token string, expiry time.Time) error {
	db := config.DB()
	revoketoken := &models.Revoke_token{
		Token:   token,
		Expires: expiry,
	}
	if err := db.Create(&revoketoken).Error; err != nil {

		return err
	}

	return nil
}
