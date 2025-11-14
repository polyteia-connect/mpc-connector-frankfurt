package measles

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
	h.router.POST("/measles-vaccination-check/schedule", h.scheduleHandler)
	h.router.GET("/measles-vaccination-check/result/:requestId", h.resultHandler)
	h.router.POST("/callback/result/:requestId", h.callbackHandler)
}

func (h *Handler) Run(addr string) error {
	slog.Info("Starting Measles API Server", "address", addr)
	return h.router.Run(addr)
}

func (h *Handler) healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *Handler) scheduleHandler(c *gin.Context) {
	payload := struct {
		RequestID    uuid.UUID   `json:"requestId"`
		FileStateIDs []uuid.UUID `json:"fileStateIds"`
	}{}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &Task{
		ID:           payload.RequestID,
		FileStateIDs: payload.FileStateIDs,
		Status:       StatusPending,
	}

	scheduledTask, err := h.service.ScheduleTask(c, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := mapTaskResult(scheduledTask)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) resultHandler(c *gin.Context) {
	requestID := c.Param("requestId")

	task, err := h.service.GetTask(c, uuid.MustParse(requestID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	response := mapTaskResult(task)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) callbackHandler(c *gin.Context) {
	requestID := uuid.MustParse(c.Param("requestId"))

	task, err := h.service.GetTask(c, requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	task.Status = StatusCompleted
	task.Result = &TaskResult{
		Match: false,
	}

	result := struct {
		Type    string `json:"type"`
		Details string `json:"details"`
	}{}

	if err := c.BindJSON(&result); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		slog.Error("Failed to bind JSON from callback", "error", err)
		return
	}

	if result.Type == "error" {
		task.Status = StatusFailed
		task.Result.Error = result.Details
		slog.Error("Task failed with error", "error", task.Result.Error, "task", task.ID)
	}

	if result.Type == "success" {
		if result.Details == "True" {
			task.Result.Match = true
		}
	}

	if err := h.service.UpdateTask(c, requestID, task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		slog.Error("Failed to update task", "error", err, "task", task.ID)
		return
	}

	slog.Info("Task completed", "task", task.ID, "result", task.Result.Match, "error", task.Result.Error)

	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

func mapTaskResult(task *Task) gin.H {
	response := gin.H{
		"status": task.Status,
	}

	if task.Result != nil {
		if task.Result.Error != "" {
			response["error"] = task.Result.Error
		}
		response["result"] = map[string]bool{
			"match": task.Result.Match,
		}
	}

	return response
}
