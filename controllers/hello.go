package controllers

import (
	"github.com/gin-gonic/gin"
)

type HelloWorldController struct {
}

func (h *HelloWorldController) Default(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hello world, climate change is real"})
}
