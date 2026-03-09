package response

import (
	 "net/http"

    "github.com/gin-gonic/gin"
)

type Response struct {
    Code    int    `json:"code"`
    Data    any    `json:"data,omitempty"`
    Message string `json:"message"`
}

type PagedData struct {
    Items any   `json:"items"`
    Total int64 `json:"total"`
    Page  int   `json:"page"`
    Size  int   `json:"size"`
}

func OK(c *gin.Context, data any) {
    c.JSON(http.StatusOK, Response{Code: http.StatusOK, Data: data, Message: "success"})
}

func Created(c *gin.Context, data any) {
    c.JSON(http.StatusCreated, Response{Code: http.StatusCreated, Data: data, Message: "created"})
}

func Paged(c *gin.Context, items any, total int64, page, size int) {
    c.JSON(http.StatusOK, Response{
        Code: http.StatusOK,
        Data: PagedData{
            Items:  items,
            Total: total,
            Page: page,
            Size: size,
        },
        Message: "success",
    })
}

func BadRequest(c *gin.Context, msg string) {
    c.JSON(http.StatusBadRequest, Response{Code: http.StatusBadRequest, Message: msg})
}

func Unauthorized(c *gin.Context) {
    c.JSON(http.StatusUnauthorized, Response{Code: http.StatusUnauthorized, Message: "unauthorized"})
}

func Forbidden(c *gin.Context) {
    c.JSON(http.StatusForbidden, Response{Code: http.StatusForbidden, Message: "forbidden"})
}
func NotFound(c *gin.Context, msg string) {
    c.JSON(http.StatusNotFound, Response{Code: http.StatusNotFound, Message: msg})
}

func Conflict(c *gin.Context, msg string) {
    c.JSON(http.StatusConflict, Response{Code: http.StatusConflict, Message: msg})
}

func InternalError(c *gin.Context) {
    c.JSON(http.StatusInternalServerError, Response{Code: http.StatusInternalServerError, Message: "internal server error"})
}

func Error(c *gin.Context, code int, msg string) {
    c.JSON(code, Response{Code: code, Message: msg})
}

func ErrorWithData(c *gin.Context, code int, msg string, data any) {
    c.JSON(code, Response{Code: code, Message: msg, Data: data})
}
