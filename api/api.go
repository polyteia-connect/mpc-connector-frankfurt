package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polyteia-de/atlas/external"
	"github.com/polyteia-de/atlas/mpc"
	"github.com/rs/xid"
)

type Handler struct {
	store          map[string]any
	engine         *gin.Engine
	externalClient *external.Client
	mpcClient      *mpc.Client
	host           string
	port           int
}

func Run(externalClient *external.Client,
	mpcClient *mpc.Client, host string, port int) {
	engine := gin.Default()

	h := &Handler{
		engine:         engine,
		externalClient: externalClient,
		mpcClient:      mpcClient,
		host:           host,
		port:           port,
		store:          make(map[string]any),
	}

	engine.Use(gin.Recovery())
	engine.POST("/launch", h.handleLaunch)
	engine.POST("/callback/:requestID", h.handleCallback)

	if err := engine.Run(fmt.Sprintf("%s:%d", h.host, h.port)); err != nil {
		panic(err)
	}
}

func (h *Handler) handleLaunch(c *gin.Context) {
	data, err := h.externalClient.FetchData(c, "/abc")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	requestID := xid.New().String()
	h.store[requestID] = data
	outputURL := fmt.Sprintf("http://%s:%d/callback/%s", h.host, h.port, requestID)

	resp, err := h.mpcClient.LaunchTask(c, data, outputURL)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) handleCallback(c *gin.Context) {
	requestID := c.Param("requestID")
	_, ok := h.store[requestID]
	if !ok {
		_ = c.AbortWithError(http.StatusNotFound, errors.New("request not found"))
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	//TODO: Do something with the MPC result
	fmt.Println(string(body))
}
