package mpc

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

//go:embed program.garble
var Program string

type LaunchRequest struct {
	ComputationID string   `json:"computation_id"`
	Participants  []string `json:"participants"`
	Program       string   `json:"program"`
	Leader        int      `json:"leader"`
	Party         int      `json:"party"`
	Input         struct {
		Array []struct {
			Array []struct {
				NumUnsigned []any `json:"NumUnsigned"`
			} `json:"Array"`
		} `json:"Array"`
	} `json:"input"`
	Output    string `json:"output,omitempty"`
	Constants struct {
		ROWS struct {
			NumUnsigned []any `json:"NumUnsigned"`
		} `json:"ROWS"`
		IDLEN struct {
			NumUnsigned []any `json:"NumUnsigned"`
		} `json:"ID_LEN"`
	} `json:"constants"`
}

type Client struct {
	client       *resty.Client
	leader       int
	party        int
	participants []string
}

func NewClient(c *resty.Client, leader, party int, participants []string) *Client {
	slog.Info("Connecting to MPC", "leader", leader, "party", party, "participants", participants)
	return &Client{
		client:       c,
		leader:       leader,
		participants: participants,
		party:        party,
	}
}

func (c *Client) LaunchTask(ctx context.Context, data []uuid.UUID, taskID uuid.UUID, callbackURL string) (string, error) {
	policy := c.createPolicy(data, taskID, callbackURL)

	resp, err := c.client.R().
		SetBody(policy).
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		Post("/schedule")
	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func (c *Client) createPolicy(data []uuid.UUID, taskID uuid.UUID, callbackURL string) LaunchRequest {
	launchRequest := LaunchRequest{}
	launchRequest.Leader = c.leader
	launchRequest.ComputationID = taskID.String()
	launchRequest.Participants = c.participants
	launchRequest.Program = Program
	launchRequest.Party = c.party
	if callbackURL != "" {
		launchRequest.Output = callbackURL
	}
	launchRequest.Constants.IDLEN.NumUnsigned = numUnsigned(16, "Usize")
	launchRequest.Constants.ROWS.NumUnsigned = numUnsigned(byte(len(data)), "Usize")

	// Initialize the outer array with the correct length
	launchRequest.Input.Array = make([]struct {
		Array []struct {
			NumUnsigned []any `json:"NumUnsigned"`
		} `json:"Array"`
	}, len(data))

	i := 0
	for _, v := range data {
		// Initialize the inner array with the correct length for HashKey
		launchRequest.Input.Array[i].Array = make([]struct {
			NumUnsigned []any `json:"NumUnsigned"`
		}, len(v))

		j := 0
		for _, bv := range v {
			launchRequest.Input.Array[i].Array[j].NumUnsigned = numUnsigned(bv, "U8")
			j++
		}
		i++
	}

	return launchRequest
}

func numUnsigned(val byte, dataType string) []any {
	return []any{
		val,
		dataType,
	}
}
