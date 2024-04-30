package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IHealthHandler interface {
	HealthCheck(c *gin.Context)
}

type healthHandler struct{}

func NewHealthHandler() IHealthHandler {
	return &healthHandler{}

}

func (h *healthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
