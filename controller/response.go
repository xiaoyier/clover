package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GlobalResponse struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
}

func Succ(ctx *gin.Context, data interface{}) {

	response(ctx, http.StatusOK, CodeSucc, "", data)
}

func Error(ctx *gin.Context, errCode ResponseCode) {

	response(ctx, http.StatusOK, errCode, "", nil)
}

func Error400(ctx *gin.Context, errCode ResponseCode) {

	response(ctx, http.StatusBadRequest, errCode, "", nil)
}

func Error401(ctx *gin.Context, errCode ResponseCode) {

	response(ctx, http.StatusUnauthorized, errCode, "", nil)
}

func Error500(ctx *gin.Context, errCode ResponseCode) {

	response(ctx, http.StatusInternalServerError, errCode, "", nil)
}

func ErrorWithMessage(ctx *gin.Context, errCode ResponseCode, message string) {

	response(ctx, http.StatusOK, errCode, message, nil)
}

func Error400WithMessage(ctx *gin.Context, errCode ResponseCode, message string) {
	response(ctx, http.StatusBadRequest, errCode, message, nil)
}

func Error401WithMessage(ctx *gin.Context, errCode ResponseCode, message string) {
	response(ctx, http.StatusUnauthorized, errCode, message, nil)
}

func Error500WithMessage(ctx *gin.Context, errCode ResponseCode, message string) {
	response(ctx, http.StatusInternalServerError, errCode, message, nil)
}

func response(ctx *gin.Context, httpCode int, errorCode ResponseCode, message string, data interface{}) {
	msg := errorCode.Message()
	if message != "" {
		msg = message
	}
	ctx.JSON(httpCode, GlobalResponse{
		Code:    errorCode,
		Message: msg,
		Data:    data,
	})
}
