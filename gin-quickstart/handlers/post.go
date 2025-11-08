package handlers

import (
	"net/http"

	"gin-quickstart/config"
	"gin-quickstart/models"
	"gin-quickstart/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreatePost 创建文章
func CreatePost(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, "用户未认证", nil)
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "请求参数解析失败", err)
		return
	}

	post.UserID = userID.(uint)

	db := config.GetDB()
	if err := db.Create(&post).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "创建文章失败", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "文章创建成功", post)
}

// GetAllPosts 获取所有文章
func GetAllPosts(c *gin.Context) {
	var posts []models.Post
	db := config.GetDB()
	if err := db.Preload("User").Find(&posts).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "获取文章列表失败", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取文章列表成功", posts)
}

// GetPost 获取单篇文章
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	db := config.GetDB()
	if err := db.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(c, http.StatusNotFound, "文章未找到", nil)
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "数据库查询失败", err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取文章成功", post)
}

// UpdatePost 更新文章
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, "用户未认证", nil)
		return
	}

	db := config.GetDB()
	var post models.Post
	if err := db.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(c, http.StatusNotFound, "文章未找到", nil)
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "数据库查询失败", err)
		}
		return
	}

	// 检查是否为作者
	if post.UserID != userID.(uint) {
		utils.HandleError(c, http.StatusForbidden, "只能更新自己的文章", nil)
		return
	}

	var input models.Post
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "请求参数解析失败", err)
		return
	}

	if err := db.Model(&post).Updates(models.Post{Title: input.Title, Content: input.Content}).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "更新文章失败", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "文章更新成功", post)
}

// DeletePost 删除文章
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		utils.HandleError(c, http.StatusUnauthorized, "用户未认证", nil)
		return
	}

	db := config.GetDB()
	var post models.Post
	if err := db.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(c, http.StatusNotFound, "文章未找到", nil)
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "数据库查询失败", err)
		}
		return
	}

	// 检查是否为作者
	if post.UserID != userID.(uint) {
		utils.HandleError(c, http.StatusForbidden, "只能删除自己的文章", nil)
		return
	}

	if err := db.Delete(&post).Error; err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "删除文章失败", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "文章删除成功", nil)
}
