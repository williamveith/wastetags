package main

import (
	"github.com/gin-gonic/gin"
)

func HomePage(c *gin.Context) (string, gin.H) {
	return "home.html", nil
}
