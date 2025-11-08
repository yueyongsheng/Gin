package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Response 统一响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应格式
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// HandleError 统一错误处理函数
func HandleError(c *gin.Context, code int, message string, err error) {
	log.Printf("Error occurred: %s | %v", message, err)

	errorResp := ErrorResponse{
		Code:    code,
		Message: message,
	}

	if err != nil {
		errorResp.Error = err.Error()
	}

	c.JSON(code, errorResp)
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}