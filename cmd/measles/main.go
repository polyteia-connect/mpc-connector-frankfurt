package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/polyteia-de/atlas/modules/measles"
	"github.com/polyteia-de/atlas/mpc"
	"github.com/polyteia-de/atlas/pkg/store"
	"github.com/polyteia-de/atlas/pkg/token"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	config := AutoConfig()
	jwtToken, err := token.NewJWT(config.JWTKeyFile, config.JWTIssuer)
	if err != nil {
		slog.Error("Failed to initialize JWT", "error", err)
		os.Exit(1)
	}
	mpcRestyClient := resty.New().
		SetBaseURL(config.MPCBaseURL).
		SetDebug(config.Debug).
		SetPreRequestHook(jwtToken.RestyMiddleware())

	taskStore := store.NewMemoryStore[*measles.Task]()

	mpcClient := mpc.NewClient(mpcRestyClient, config.MPCLeaderID, config.MPCPartyID, config.MPCParticipants)

	measlesService := measles.NewService(taskStore, mpcClient, config.CallbackBaseURL, config.ESUBaseURL)

	measlesHandler := measles.NewHandler(measlesService)
	measlesHandler.Configure()

	if err := measlesHandler.Run(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
		slog.Error("Failed to run Measles handler", "error", err)
		os.Exit(1)
	}
}
