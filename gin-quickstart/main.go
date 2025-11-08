package main

import (
	"fmt"
	"log"

	"gin-quickstart/config"
	"gin-quickstart/handlers"
	"gin-quickstart/middleware"
	"gin-quickstart/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	config.InitDB()

	// 自动迁移模型
	db := config.GetDB()
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	log.Println("Database migrated successfully")

	// 创建 Gin 路由
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())

	// 公开路由
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/posts", handlers.GetAllPosts)
	r.GET("/posts/:id", handlers.GetPost)
	r.GET("/comments/post/:post_id", handlers.GetPostComments)

	// 需要认证的路由
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/create-post", handlers.CreatePost)
		auth.PUT("/upposts/:id", handlers.UpdatePost)
		auth.DELETE("/delposts/:id", handlers.DeletePost)
		auth.POST("/comments", handlers.CreateComment)
	}

	// 记录服务启动信息
	fmt.Println("Server starting on :8080")
	log.Println("服务已启动,监听端口 :8080")

	// 启动服务器
	log.Println("Server starting on http://localhost:8080")
	r.Run(":8080")
}
