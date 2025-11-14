package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/polyteia-de/atlas/modules/esu"
	"github.com/polyteia-de/atlas/mpc"
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
	vaccinationRestyClient := resty.New().
		SetBaseURL(config.VaccinationBaseURL).
		SetDebug(config.Debug).
		SetPreRequestHook(jwtToken.RestyMiddleware())

	esuClient := esu.NewClient(vaccinationRestyClient)

	mpcRestyClient := resty.New().
		SetBaseURL(config.MPCBaseURL).
		SetDebug(config.Debug).
		SetPreRequestHook(jwtToken.RestyMiddleware())

	mpcClient := mpc.NewClient(mpcRestyClient, config.MPCLeaderID, config.MPCPartyID, config.MPCParticipants)

	esuService := esu.NewService(mpcClient, esuClient)

	esuHandler := esu.NewHandler(esuService)
	esuHandler.Configure()

	if err := esuHandler.Run(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
		slog.Error("Failed to run ESU handler", "error", err)
		os.Exit(1)
	}
}
