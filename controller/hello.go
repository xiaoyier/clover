package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HelloHandler(c *gin.Context) {

	c.String(http.StatusOK, "hello world!")
}
