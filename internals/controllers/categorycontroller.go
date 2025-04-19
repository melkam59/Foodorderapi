package controllers

import (
	"foodorderapi/internals/config"
	"foodorderapi/internals/models"
	"net/http"

	"github.com/labstack/echo/v4"
)


func CreateCategory(c echo.Context)error{

	role  := c.Get("role").(string)
	merchantID := c.Get("merchantID").(string)
	db:=config.DB()

	var merchants *models.Merchant
	var category *models.Category
	exitingcategory :=&models.Category{}

	if role!="merchant"{
		data:=map[string]interface{}{
		"message": "Access denied. Only merchants can perform this operation.",
		}

		return c.JSON(http.StatusForbidden ,data)
	}

		if res := db.Where("id = ?", merchantID).Find(&merchants); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if err := c.Bind(&category);err!=nil{
		return c.String(http.StatusBadRequest,"invalid request payload")
	}


	if res:=db.Where("categoryname=? AND merchant_id=?", category.Categoryname,merchantID).First(&exitingcategory);res.Error ==nil{
		data := map[string]interface{}{
			"message": "Category already exists for the provided food name and merchant ID",
		}
		return c.JSON(http.StatusConflict, data)
	}

	newcategory :=&models.Category{

	Categoryname :category.Categoryname,
	Categorydescription :category.Categorydescription,
	Categoryimage :category.Categoryimage,
	MerchantID        :merchantID,
	MerchantShortCode :merchants.MerchantShortcode,

	}
	if err:=db.Create(&newcategory).Error;err !=nil{

	return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, newcategory)

}



func GetCategory(c echo.Context) error {
	role := c.Get("role").(string)
	merchantID := c.Get("merchantID").(string)
	var merchants models.Merchant
	var categories []models.Category
	var menus []models.Menu // New variable to hold the menus

	db := config.DB()

	if role != "merchant" {
		data := map[string]interface{}{
			"message": "Access denied. Only merchants can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	if res := db.Where("id = ?", merchantID).First(&merchants); res.Error != nil {
		data := map[string]interface{}{
			"message": "Merchant not found",
		}
		return c.JSON(http.StatusNotFound, data)
	}

	if res := db.Where("merchant_id = ?", merchantID).Find(&categories); res.Error != nil {
		data := map[string]interface{}{
			"message": "Category not found",
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Retrieve menus for each category
for i := range categories {
	category := &categories[i]
	if res := db.Where("food_category = ? AND merchant_id = ?", category.Categoryname, merchantID).Find(&menus); res.Error != nil {
		data := map[string]interface{}{
			"message": "Menus not found for category",
		}
		return c.JSON(http.StatusNotFound, data)
	}
	category.Menu = menus // Assign menus to the category
	menus = []models.Menu{} // Reset the menus slice for the next category
}

	return c.JSON(http.StatusOK, categories)
}




func EditCategory(c echo.Context)error{

	
	role:=c.Get("role").(string)
	merchantid  :=c.Get("merchantID").(string)
	var merchants *models.Merchant
	var category  models.Category
    var payload models.UpdateCategory
	categoryID:=c.Param("id")

	db:=config.DB()


    if role != "merchant" {
		data := map[string]interface{}{
			"message": "Access denied. Only merchants can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	if res:=db.Where("id =?",merchantid).Find(&merchants);res.Error!=nil{
		data:=map[string]interface{}{
			"message":"merchant not found",
		}

		return c.JSON(http.StatusNotFound,data)
	}

	if res := db.Where("id = ? AND merchant_id = ?", categoryID, merchantid).First(&category); res.Error != nil {
		data := map[string]interface{}{
			"message": "Category not found",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	if payload.Categoryname != "" {
		category.Categoryname = payload.Categoryname
		if err := db.Model(&models.Menu{}).Where("food_category = ?", category.Categoryname).Update("food_category", payload.Categoryname).Error; err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update food_category in menu table.")
		}
	}
	if payload.Categorydescription != "" {
		category.Categorydescription=payload.Categorydescription
	}
	if payload.Categoryimage != ""{
		category.Categoryimage=payload.Categoryimage
	}
	
	if err := db.Save(&category).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Failed to update category. Please try again with a new category name.")
	}





	

	return c.JSON(http.StatusOK, category)



}


func DeleteCategory(c echo.Context) error {
	role := c.Get("role").(string)

	// Check if the role is merchant
	if role != "merchant" {
		data := map[string]interface{}{
			"message": "Access denied. Only merchants can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	db := config.DB()

	merchantID := c.Get("merchantID").(string)
	categoryID := c.Param("id")

	var category models.Category

	// Find the menu item by ID and associated merchant ID
	if res := db.Where("id = ? AND merchant_id = ?", categoryID, merchantID).First(&category); res.Error != nil {
		data := map[string]interface{}{
			"message": "category item not found",
		}
		return c.JSON(http.StatusNotFound, data)
	}



	if res := db.Where("food_category = ?", category.Categoryname).Delete(&models.Menu{}); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}
	// Delete the menu item
	if res := db.Delete(&category); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	data := map[string]interface{}{
		"message": "Category item deleted successfully",
	}
	return c.JSON(http.StatusOK, data)
}



//for users 

func DisplayCategory(c echo.Context)error{


	db:=config.DB()
	var category []models.Category
	var merchant models.Merchant

var reqBody struct {
		MerchantShortcode int64 `json:"merchantshortcode"`
	}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}


if res:=db.Where("merchant_shortcode=?",reqBody.MerchantShortcode).Find(&merchant);  res.Error!=nil{
data:=map[string]interface{}{
		"message":"merchant not found",
	}
	return c.JSON(http.StatusInternalServerError,data)
}

if res:=db.Where("merchant_short_code=?",reqBody.MerchantShortcode).Find(&category);res.Error!=nil{
	data:=map[string]interface{}{
		"message":"category not found for this merchant",
	}
	return c.JSON(http.StatusInternalServerError,data)
}


  




return  c.JSON(http.StatusOK,category)




}

func FoodNumberByCategory(c echo.Context) error {
	db := config.DB()
	var menu []models.Menu
	var categories []models.Category

	var reqBody struct {
		MerchantShortcode int64 `json:"merchantshortcode"`
	}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if res := db.Where("merchant_short_code = ?", reqBody.MerchantShortcode).Find(&categories); res.Error != nil {
		data := map[string]interface{}{
			"message": "Categories not found for this merchant",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	menuCountByCategory := make(map[string]int)

	for _, category := range categories {
		if res := db.Where("merchant_short_code = ? AND food_category = ?", reqBody.MerchantShortcode, category.Categoryname).Find(&menu); res.Error != nil {
			data := map[string]interface{}{
				"message": "Menus not found for this category",
			}
			return c.JSON(http.StatusInternalServerError, data)
		}

		menuCountByCategory[category.Categoryname] = len(menu)
	}

	return c.JSON(http.StatusOK, menuCountByCategory)
}


func NumberofCategoriesforMerchant(c echo.Context)  error{
	db:=config.DB()
    var category []models.Category
    role := c.Get("role").(string)

	// Check if the role is merchant
	if role != "merchant" {
		data := map[string]interface{}{
			"message": "Access denied. Only merchants can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	

	merchantID := c.Get("merchantID").(string)


	if res := db.Where("merchant_id = ?",merchantID).Find(&category); res.Error != nil {
		data := map[string]interface{}{
			"message": "Categories not found for this merchant",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

   categoryCount := len(category)

	// data := map[string]interface{}{
	// 	"message":         "Success",
	// 	"categoryCount":   categoryCount,
	// 	"merchantShortcode": reqBody.MerchantShortcode,
	// }

	return c.JSON(http.StatusOK, categoryCount)


}



func NumberofCategories(c echo.Context) error {
	db := config.DB()

	var category []models.Category
	var reqBody struct {
		MerchantShortcode int64 `json:"merchantshortcode"`
	}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if res := db.Where("merchant_short_code = ?", reqBody.MerchantShortcode).Find(&category); res.Error != nil {
		data := map[string]interface{}{
			"message": "Categories not found for this merchant",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	categoryCount := len(category)

	// data := map[string]interface{}{
	// 	"message":         "Success",
	// 	"categoryCount":   categoryCount,
	// 	"merchantShortcode": reqBody.MerchantShortcode,
	// }

	return c.JSON(http.StatusOK, categoryCount)
}


func NumberofMenusforMerchant(c echo.Context)  error{
	db:=config.DB()
    var menu []models.Menu
    role := c.Get("role").(string)

	// Check if the role is merchant
	if role != "merchant" {
		data := map[string]interface{}{
			"message": "Access denied. Only merchants can perform this operation.",
		}
		return c.JSON(http.StatusForbidden, data)
	}

	

	merchantID := c.Get("merchantID").(string)


	if res := db.Where("merchant_id = ?",merchantID).Find(&menu); res.Error != nil {
		data := map[string]interface{}{
			"message": "Categories not found for this merchant",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

   menuCount := len(menu)

	// data := map[string]interface{}{
	// 	"message":         "Success",
	// 	"categoryCount":   categoryCount,
	// 	"merchantShortcode": reqBody.MerchantShortcode,
	// }

	return c.JSON(http.StatusOK, menuCount)


}


func NumberofMenus(c echo.Context) error {
	db := config.DB()

	var menu []models.Menu
	var reqBody struct {
		MerchantShortcode int64 `json:"merchantshortcode"`
	}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if res := db.Where("merchant_short_code = ?", reqBody.MerchantShortcode).Find(&menu); res.Error != nil {
		data := map[string]interface{}{
			"message": "Categories not found for this merchant",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	menuCount := len(menu)

	// data := map[string]interface{}{
	// 	"message":         "Success",
	// 	"menuCount":   menuCount,
	// 	"merchantShortcode": reqBody.MerchantShortcode,
	// }

	return c.JSON(http.StatusOK, menuCount)
}