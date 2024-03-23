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

func CreateComment(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)

	comment := models.Comment{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&comment)
	} else {
		ctx.ShouldBind(&comment)
	}

	comment.UserID = userID

	err := db.Debug().Create(&comment).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         comment.ID,
		"message":    comment.Message,
		"photo_id":   comment.PhotoID,
		"user_id":    comment.UserID,
		"created_at": comment.CreatedAt,
	})
}

func GetComments(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	photo := models.Photo{}
	user := models.User{}
	comment := []models.Comment{}

	userID := uint(userData["id"].(float64))

	err := db.Where("user_id = ?", userID).Find(&comment).Error
	errUser := db.Where("id = ?", userID).Find(&user).Error
	errPhoto := db.Where("user_id = ?", userID).Find(&photo).Error

	if err != nil && errUser != nil && errPhoto != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	for _, value := range comment {
		ctx.JSON(http.StatusOK, gin.H{
			"id":         value.ID,
			"message":    value.Message,
			"photo_id":   value.PhotoID,
			"user_id":    value.UserID,
			"updated_at": value.UpdatedAt,
			"created_at": value.CreatedAt,
			"User": gin.H{
				"id":       user.ID,
				"email":    user.Email,
				"username": user.Username,
			},
			"Photo": gin.H{
				"id":        photo.ID,
				"title":     photo.Title,
				"caption":   photo.Caption,
				"photo_url": photo.PhotoURL,
				"user_id":   photo.UserID,
			},
		})
	}

}

func UpdateComment(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(ctx)
	comment := models.Comment{}

	commentId, _ := strconv.Atoi(ctx.Param("commentId"))
	userId := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&comment)
	} else {
		ctx.ShouldBind(&comment)
	}

	comment.ID = uint(commentId)
	comment.UserID = uint(userId)

	err := db.Model(&comment).Where("id = ?", commentId).Updates(models.Comment{Message: comment.Message}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

func DeleteComment(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	comment := models.Comment{}

	commentId, _ := strconv.Atoi(ctx.Param("commentId"))
	userId := uint(userData["id"].(float64))

	comment.ID = uint(commentId)
	comment.UserID = userId

	err := db.Model(&comment).Where("id = ?", commentId).Delete(models.Comment{}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
