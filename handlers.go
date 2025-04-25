package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MetadataHandler(metadata Metadata) func(c *gin.Context) {
	return func(c *gin.Context) {

		c.JSON(http.StatusOK, metadata)
	}
}
