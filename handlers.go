package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func reverseStrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Param("str")
		if s == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no string provided"})
			return
		}

		rev, _ := reverseStrSvc(c.Request.Context(), s)

		c.JSON(http.StatusOK, gin.H{"string": rev})
	}
}
