package esu

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *Handler) Run(addr string) error {
	slog.Info("Starting ESU API Server", "address", addr)
	return h.router.Run(addr)
}

func (h *Handler) healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *Handler) scheduleHandler(c *gin.Context) {
	taskID := uuid.MustParse(c.Param("taskId"))

	err := h.service.ScheduleTask(c, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "scheduled"})
}
