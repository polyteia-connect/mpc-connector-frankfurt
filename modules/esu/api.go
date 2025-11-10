package esu

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	router  *gin.Engine
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
		router:  gin.Default(),
	}
}

func (h *Handler) Configure() {
	h.router.Use(gin.Recovery(), gin.Logger())
	h.router.GET("/health", h.healthHandler)
	h.router.PUT("/schedule/:taskId", h.scheduleHandler)
}

func (h *Handler) healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *Handler) scheduleHandler(c *gin.Context) {
	taskID := c.Param("taskId")

	err := h.service.ScheduleTask(c, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "scheduled"})
}
