package esu

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewClient(c *resty.Client) *Client {
	return &Client{
		client: c,
	}
}

func (c *Client) GetVaccinatedIDs(ctx context.Context) ([]string, error) {
	resp, err := c.client.R().SetContext(ctx).SetResult(&[]string{}).Get("/vaccinated-ids")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get vaccinated IDs: %d, %s", resp.StatusCode(), resp.String())
	}

	ids, ok := resp.Result().(*[]string)
	if !ok {
		return nil, fmt.Errorf("invalid response type: %T", resp.Result())
	}

	return *ids, nil
}
