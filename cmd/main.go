package main

import (
	"github.com/polyteia-de/atlas/api"
	"github.com/polyteia-de/atlas/config"
	"github.com/polyteia-de/atlas/external"
	"github.com/polyteia-de/atlas/mpc"
	"resty.dev/v3"
)

func main() {
	cfg := config.Auto()
	externalClientResty := resty.New().SetBaseURL(cfg.ExternalClientURL).SetDebug(cfg.Debug)
	externalClient := external.NewClient(externalClientResty)

	mpcClientResty := resty.New().SetBaseURL(cfg.MPCBaseURL).SetDebug(cfg.Debug)
	mpcClient := mpc.NewClient(mpcClientResty, cfg.MPCLeaderID, cfg.MPCPartyID, cfg.MPCParticipants)

	api.Run(externalClient, mpcClient, cfg.Host, cfg.Port)
}
