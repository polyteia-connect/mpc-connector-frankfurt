package measles

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/polyteia-de/atlas/mpc"
	"github.com/polyteia-de/atlas/pkg/store"
)

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

type TaskResult struct {
	Match bool
	Error string
}

type Task struct {
	ID           uuid.UUID
	FileStateIDs []uuid.UUID
	Status       Status
	Result       *TaskResult
}

type Service struct {
	store           store.Store[*Task]
	mpcClient       *mpc.Client
	callbackBaseURL string
	esuClient       *resty.Client
}

func NewService(store store.Store[*Task], mpcClient *mpc.Client, callbackBaseURL string, esuBaseURL string) *Service {
	esuClient := resty.New().SetBaseURL(esuBaseURL)
	return &Service{
		store:           store,
		mpcClient:       mpcClient,
		callbackBaseURL: callbackBaseURL,
		esuClient:       esuClient,
	}
}

func (s *Service) GetTask(ctx context.Context, requestID uuid.UUID) (*Task, error) {
	return s.store.Get(ctx, requestID.String())
}

func (s *Service) UpdateTask(ctx context.Context, requestID uuid.UUID, task *Task) error {
	return s.store.Set(ctx, requestID.String(), task)
}

func (s *Service) ScheduleTask(ctx context.Context, task *Task) (*Task, error) {
	// Check if task is already there with the given ID.
	t, err := s.store.Get(ctx, task.ID.String())
	if err != nil {
		return nil, err
	}
	if t != nil {
		return t, nil
	}

	// Store the task
	if err := s.store.Set(ctx, task.ID.String(), task); err != nil {
		return nil, err
	}

	// Launch the MPC task in the background
	go s.launchMPCTask(ctx, task)

	// Launch the ESU task in the background
	go s.launchESUTask(ctx, task)

	return task, nil
}

func (s *Service) launchESUTask(ctx context.Context, task *Task) {
	path := fmt.Sprintf("/schedule/%s", task.ID.String())
	resp, err := s.esuClient.R().SetContext(ctx).Put(path)
	if err != nil {
		slog.Error("Failed to launch ESU task", "error", err, "task", task.ID)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error("Failed to launch ESU task", "error", resp.String(), "task", task.ID)
		return
	}

	slog.Info("ESU task launched", "task", task.ID)
}

func (s *Service) launchMPCTask(ctx context.Context, task *Task) {
	callbackURL := fmt.Sprintf("%s/callback/result/%s", s.callbackBaseURL, task.ID.String())

	_, err := s.mpcClient.LaunchTask(ctx, task.FileStateIDs, task.ID, callbackURL)
	if err != nil {
		slog.Error("Failed to launch MPC task", "error", err, "task", task.ID)
		return
	}

	slog.Info("MPC task launched", "task", task.ID)
}
