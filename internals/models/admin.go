package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	Id          string `gorm:"primary_key;" json:"id"`
	AdminName   string ` validate:" required , max=30"  json:"name" `
	Password    string ` validate:" required , max=30, min=6 "  json:"password" `
	Phonenumber int64  `gorm:"uniqueIndex" validate:"required,ethiopianPhoneNumber" json:"phonenumber"`
	Email       string ` gorm:"uniqueIndex;not null" json:"email,omitempty"`
	PrivateKey  string
	PublicKey   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (base *Admin) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	base.Id = string(uuid)
	return nil
}

type General struct {
	Id                string `gorm:"primary_key;" json:"id"`
	MerchantID        string `json:"merchantid" gorm:"foreignkey"`
	AdminID           string `json:"adminid" gorm:"foreignkey"`
	Name              string ` validate:" required , max=30"  json:"name" `
	Password          string ` validate:" required , max=30, min=6 "  json:"password" `
	Phonenumber       int64  `gorm:"uniqueIndex" validate:"required,ethiopianPhoneNumber" json:"phonenumber"`
	Email             string ` gorm:"uniqueIndex;not null" json:"email,omitempty"`
	Role              string ` validate:" required , max=10"  json:"role" `
	MerchantShortcode int64  `gorm:"unique_index;foreignkey:MerchantID" json:"merchantshortcode"`
	IsActive          bool   `json:"isActive"`
	PrivateKey        string
	PublicKey         string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (base *General) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	base.Id = string(uuid)
	return nil
}

type Revoke_token struct {
	Id      string
	Token   string
	Expires time.Time
}

func (base *Revoke_token) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	base.Id = string(uuid)
	return nil
}
