package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Controller) RenderIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Hub Dashboard",
	})
}
