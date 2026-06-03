package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Code: 0, Message: "created", Data: data})
}

func Fail(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, Response{Code: -1, Message: msg})
}

func BadRequest(c *gin.Context, msg string) {
	Fail(c, http.StatusBadRequest, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Fail(c, http.StatusUnauthorized, msg)
}

func Forbidden(c *gin.Context, msg string) {
	Fail(c, http.StatusForbidden, msg)
}

func NotFound(c *gin.Context, msg string) {
	Fail(c, http.StatusNotFound, msg)
}

func InternalError(c *gin.Context, msg string) {
	Fail(c, http.StatusInternalServerError, msg)
}
