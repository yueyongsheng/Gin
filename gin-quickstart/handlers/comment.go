package handlers

import (
	"net/http"

	"gin-quickstart/config"
	"gin-quickstart/models"
	"gin-quickstart/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, "用户未认证", nil)
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "请求参数解析失败", err)
		return
	}

	db := config.GetDB()
	// 检查文章是否存在
	var post models.Post
	if err := db.First(&post, comment.PostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(c, http.StatusNotFound, "文章未找到", nil)
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "数据库查询失败", err)
		}
		return
	}

	comment.UserID = userID.(uint)

	if err := db.Create(&comment).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "创建评论失败", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "评论创建成功", comment)
}

// GetPostComments 获取文章的所有评论
func GetPostComments(c *gin.Context) {
	postID := c.Param("post_id")
	var comments []models.Comment

	db := config.GetDB()
	// 检查文章是否存在
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(c, http.StatusNotFound, "文章未找到", nil)
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "数据库查询失败", err)
		}
		return
	}

	if err := db.Where("post_id = ?", postID).Preload("User").Find(&comments).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "获取评论列表失败", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取评论列表成功", comments)
}
