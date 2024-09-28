package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        string    `gorm:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
}

func (base *Model) BeforeCreate(scope *gorm.DB) error {
	uuid, err := uuid.New().MarshalText()
	if err != nil {
		return err
	}
	base.ID = string(uuid)
	return nil
}

// Menu represents a food menu item
type Menu struct {
	Model
	FoodName          string  `json:"name"`
	Ingredients       string  `json:"ingredients"`
	Price             float64 `json:"price"`
	Image             string  `json:"image"`
	MerchantID        string  `json:"merchantid" gorm:"foreignkey"`
	MerchantShortCode int64   ` gorm:"foreignkey" json:"merchantshortcode"`
    FoodCategory       string ` json:"foodcategory"`
	IsFasting          bool    `json:"is_fasting"`
}

type UpdateMenu struct {
	FoodName    string  `json:"name"`
	Ingredients string  `json:"ingredients"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
	FoodCategory       string ` json:"foodcategory"`
}

// Order represents a food order
type Order struct {
	Model
	MenuID    string  `json:"menu_id"`
	Quantity  int     `json:"quantity"`
	TotalCost float64 `json:"total_cost"`
}

// Customer represents a customer who places an order
type Customer struct {
	Model
	Name        string  `gorm:"size:30;not null" json:"Name"`
	PhoneNumber string  `gorm:"size:13;not null" json:"phoneNumber"`
	Email       *string `gorm:"size:50" json:"email,omitempty"`
	Address     string  `gorm:"size:255;" json:"address"`
}



type Category struct{
	Model
	Categoryname string `json:"categoryname" gorm:"not null"`
	Categorydescription string `json:"categorydescription"`
	Categoryimage string `json:"categoryimage"`
	MerchantID        string  `json:"merchantid" gorm:"foreignkey"`
	MerchantShortCode int64   ` gorm:"foreignkey" json:"merchantshortcode"`
	Menu                []Menu `gorm:"-" json:"menu"`
}

type UpdateCategory struct{
	
	Categoryname string `json:"categoryname" gorm:"not null"`
	Categorydescription string `json:"categorydescription"`
	Categoryimage string `json:"categoryimage"`
	
}















// OrderRequest represents the request payload for creating an order
type OrderRequest struct {
	MenuID   string `json:"menu_id"`
	Quantity int    `json:"quantity"`
	Customer struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Phone   string `json:"phone"`
	} `json:"customer"`
}


