package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController is controller for health check
type HealthController struct{}

// Status returns Working!
func (h HealthController) Status(c *gin.Context) {
	c.String(http.StatusOK, "Working!")
}
