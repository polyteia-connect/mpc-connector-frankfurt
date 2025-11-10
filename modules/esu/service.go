package esu

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/polyteia-de/atlas/mpc"
)

type Service struct {
	mpcClient       *mpc.Client
	esuClient       *Client
	callbackBaseURL string
}

func NewService(mpcClient *mpc.Client, esuClient *Client, callbackBaseURL string) *Service {
	return &Service{
		mpcClient:       mpcClient,
		esuClient:       esuClient,
		callbackBaseURL: callbackBaseURL,
	}
}

func (s *Service) ScheduleTask(ctx context.Context, taskID string) error {
	ids, err := s.esuClient.GetVaccinatedIDs(ctx)
	if err != nil {
		return err
	}

	s.launchMPCTask(ctx, taskID, ids)

	return nil
}

func (s *Service) launchMPCTask(ctx context.Context, taskID string, ids []string) {
	callbackURL := fmt.Sprintf("%s/callback/result/%s", s.callbackBaseURL, taskID)
	data := make([][]byte, len(ids))

	for i, id := range ids {
		data[i] = []byte(id)
	}

	_, err := s.mpcClient.LaunchTask(ctx, data, callbackURL)
	if err != nil {
		slog.Error("Failed to launch MPC task", "error", err, "task", taskID)
		return
	}

	slog.Info("MPC task launched", "task", taskID)
}
