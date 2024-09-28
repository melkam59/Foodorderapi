package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Merchant struct {
	Id                string `gorm:"primary_key;" json:"id"`
	BusinessName      string ` validate:" required , max=30"  json:"businessname" `
	OwnerName         string ` validate:" required , max=30"  json:"ownername" `
	ContactPerson     string ` validate:" required , max=30"  json:"contactperson" `
	Password          string ` validate:" required , max=30, min=6 "  json:"password" `
	Phonenumber       int64  `gorm:"uniqueIndex" validate:"required,ethiopianPhoneNumber" json:"phonenumber"`
	Token             string `json:"token"`
	Email             string ` validate:" required" gorm:"uniqueIndex;not null" json:"email,omitempty"`
	MerchantShortcode int64  `gorm:"unique_index" validate:"required" json:"merchantshortcode"`
	IsActive          bool   `json:"isActive"`
	IsUpdated         bool   `json:"isUpdated"`

	PrivateKey string
	PublicKey  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (base *Merchant) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	base.Id = string(uuid)
	return nil
}

type Signupinputs struct {
	BusinessName      string ` validate:" required , max=30"  json:"businessname" `
	OwnerName         string ` validate:" required , max=30"  json:"ownername" `
	ContactPerson     string ` validate:" required , max=30"  json:"contactperson" `
	Email             string `validate:" required" gorm:"uniqueIndex;not null" json:"email"`
	Password          string ` validate:" required , max=30, min=6 "  json:"password" `
	Phonenumber       int64  `gorm:"uniqueIndex" validate:"required,ethiopianPhoneNumber" json:"phonenumber"`
	MerchantShortcode int64  `gorm:"unique_index" validate:"required" json:"merchantshortcode"`
	PrivateKey        string
	PublicKey         string
}

type Signininputone struct {
	Phonenumber int64  `gorm:"uniqueIndex" validate:"required,ethiopianPhoneNumber" json:"phonenumber"`
	Password    string ` validate:" required , max=30, min=6 "  json:"password" `
}

type Merchantsignin struct {
	MerchantShortcode int64  `gorm:"unique_index" validate:"required" json:"merchantshortcode"`
	Password          string ` validate:" required , max=30, min=6 "  json:"password"`
}

type Signininputs struct {
	MerchantShortcode int64  `gorm:"unique_index" validate:"required" json:"merchantshortcode"`
}

type MerchantResponse struct {
	Id                string `json:"id,omitempty"`
	BusinessName      string ` validate:" required , max=30"  json:"businessname" `
	OwnerName         string ` validate:" required , max=30"  json:"ownername" `
	ContactPerson     string ` validate:" required , max=30"  json:"contactperson" `
	Email             string `validate:" required" gorm:"uniqueIndex;not null" json:"email"`
	Phonenumber       int64  `gorm:"uniqueIndex" validate:"required,ethiopianPhoneNumber" json:"phonenumber"`
	IsUpdated         bool   `json:"isUpdated"`
	PublicKey         string
	MerchantShortcode int64     `gorm:"unique_index" validate:"required" json:"merchantshortcode"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
