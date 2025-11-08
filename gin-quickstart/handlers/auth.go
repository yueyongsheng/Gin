package handlers

import (
	"net/http"
	"strings"
	"time"

	"gin-quickstart/config"
	"gin-quickstart/middleware"
	"gin-quickstart/models"
	"gin-quickstart/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register 用户注册
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "请求参数解析失败", err)
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "密码加密失败", err)
		return
	}
	user.Password = string(hashedPassword)

	db := config.GetDB()
	if err := db.Create(&user).Error; err != nil {
		// 检查是否为唯一性约束违反错误
		if strings.Contains(err.Error(), "Duplicate entry") {
			utils.HandleError(c, http.StatusConflict, "用户名或邮箱已存在", err)
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "创建用户失败", err)
		}
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "用户注册成功", nil)
}

// Login 用户登录
func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "请求参数解析失败", err)
		return
	}

	db := config.GetDB()
	var storedUser models.User
	if err := db.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.HandleError(c, http.StatusUnauthorized, "用户名或密码错误", nil)
		} else {
			utils.HandleError(c, http.StatusInternalServerError, "数据库查询失败", err)
		}
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "用户名或密码错误", nil)
		return
	}

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       storedUser.ID,
		"username": storedUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(middleware.SECRET_KEY))
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "生成令牌失败", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "登录成功", map[string]interface{}{
		"token":   tokenString,
		"user_id": storedUser.ID,
	})
}
