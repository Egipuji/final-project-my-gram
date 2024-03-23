package controllers

import (
	"my-gram/database"
	"my-gram/helpers"
	"my-gram/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AddPhoto(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)

	photo := models.Photo{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&photo)
	} else {
		ctx.ShouldBind(&photo)
	}

	photo.UserID = userID

	err := db.Debug().Create(&photo).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoURL,
		"user_id":    photo.UserID,
		"created_at": photo.CreatedAt,
	})
}

func GetPhotos(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	photo := []models.Photo{}
	user := models.User{}

	userID := uint(userData["id"].(float64))

	err := db.Where("user_id = ?", userID).Find(&photo).Error
	errUser := db.Where("id = ?", userID).Find(&user).Error

	if err != nil || errUser != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	if len(photo) <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "anda tidak ada foto",
		})
		return
	}

	for _, value := range photo {
		ctx.JSON(http.StatusOK, gin.H{
			"id":         value.ID,
			"title":      value.Title,
			"caption":    value.Caption,
			"photo_url":  value.PhotoURL,
			"user_id":    value.UserID,
			"created_at": value.CreatedAt,
			"updated_at": value.UpdatedAt,
			"User": gin.H{
				"email":    user.Email,
				"username": user.Username,
			},
		})
	}

}

func UpdatePhoto(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)
	photo := models.Photo{}

	photoId, _ := strconv.Atoi(ctx.Param("photoId"))
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&photo)
	} else {
		ctx.ShouldBind(&photo)
	}

	photo.UserID = userId
	photo.ID = uint(photoId)

	err := db.Model(&photo).Where("id = ?", photoId).Updates(models.Photo{Title: photo.Title, Caption: photo.Caption, PhotoURL: photo.PhotoURL}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoURL,
		"user_id":    photo.UserID,
		"updated_at": photo.UpdatedAt,
	})
}

func DeletePhoto(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	photo := models.Photo{}

	photoId, _ := strconv.Atoi(ctx.Param("photoId"))
	userId := uint(userData["id"].(float64))

	photo.ID = uint(photoId)
	photo.UserID = userId

	err := db.Model(&photo).Where("id = ?", photoId).Delete(models.Photo{}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}
