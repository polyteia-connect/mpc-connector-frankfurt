package measles

import (
	"fmt"
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
	h.router.POST("/measles-vaccination-check/schedule", h.scheduleHandler)
	h.router.GET("/measles-vaccination-check/result/:requestId", h.resultHandler)
	h.router.GET("/callback/result/:requestId", h.callbackHandler)
}

func (h *Handler) healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *Handler) scheduleHandler(c *gin.Context) {
	payload := struct {
		RequestID    string   `json:"requestId"`
		FileStateIDs []string `json:"fileStateIds"`
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

	task, err := h.service.GetTask(c, requestID)
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
	requestID := c.Param("requestId")

	task, err := h.service.GetTask(c, requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("callback", string(body))

	task.Status = StatusCompleted
	h.service.UpdateTask(c, requestID, task)

	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

func mapTaskResult(task *Task) gin.H {
	response := gin.H{
		"status": task.Status,
	}

	if task.Result != nil {
		response["result"] = map[string]bool{
			"match": task.Result.Match,
		}
	}

	return response
}
