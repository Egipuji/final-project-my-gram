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

func CreateSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)

	SocialMedia := models.SocialMedia{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&SocialMedia)
	} else {
		ctx.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserID = userID

	err := db.Debug().Create(&SocialMedia).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, SocialMedia)
}

func GetSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	socialMedia := models.SocialMedia{}
	user := models.User{}

	userID := uint(userData["id"].(float64))

	socialMedia.UserID = userID

	err := db.Where("user_id = ?", userID).Find(&socialMedia).Error
	errUser := db.Where("id = ?", userID).Find(&user).Error

	if err != nil && errUser != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"social_media": gin.H{
			"id":               socialMedia.ID,
			"name":             socialMedia.Name,
			"social_media_url": socialMedia.SocialMediaURL,
			"user_id":          socialMedia.UserID,
			"created_at":       socialMedia.CreatedAt,
			"updated_at":       socialMedia.UpdatedAt,
			"User": gin.H{
				"id":       user.ID,
				"username": user.Username,
			},
		},
	})
}

func UpdateSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)
	socialMedia := models.SocialMedia{}

	socialMediaID, _ := strconv.Atoi(ctx.Param("socialMediaId"))
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&socialMedia)
	} else {
		ctx.ShouldBind(&socialMedia)
	}

	socialMedia.UserID = userId
	socialMedia.ID = uint(socialMediaID)

	err := db.Model(&socialMedia).Where("id = ?", socialMediaID).Updates(models.SocialMedia{Name: socialMedia.Name, SocialMediaURL: socialMedia.SocialMediaURL}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"updated_at":       socialMedia.UpdatedAt,
	})
}

func DeleteSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	socialMedia := models.SocialMedia{}

	socialMediaId, _ := strconv.Atoi(ctx.Param("socialMediaId"))
	userId := uint(userData["id"].(float64))

	socialMedia.ID = uint(socialMediaId)
	socialMedia.UserID = userId

	err := db.Model(&socialMedia).Where("id = ?", socialMediaId).Delete(models.SocialMedia{}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
