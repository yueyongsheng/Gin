package models

import "gorm.io/gorm"

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}

// Post 文章模型
type Post struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Content  string `gorm:"not null"`
	UserID   uint
	User     User
	Comments []Comment // 一对多：文章有多个评论
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}
