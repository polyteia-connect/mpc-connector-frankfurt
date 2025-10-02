package external

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"resty.dev/v3"
)

type Data struct {
	HashKey     []byte `json:"-"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
}

type Client struct {
	client *resty.Client
}

func NewClient(c *resty.Client) *Client {
	return &Client{
		client: c,
	}
}

func (c *Client) FetchData(ctx context.Context, path string) (map[string]Data, error) {
	// TODO: Call actual endpoint
	data := make(map[string]Data)
	for i := range 100 {
		tuple := Data{
			FirstName:   fmt.Sprintf("First name %d", i),
			LastName:    fmt.Sprintf("Last name %d", i),
			DateOfBirth: fmt.Sprintf("Date of birth %d", i),
		}
		hashKey := shaHash(tuple)
		tuple.HashKey = hashKey
		data[hex.EncodeToString(hashKey)] = tuple
	}

	return data, nil
}

func shaHash(data Data) []byte {
	hash := sha256.New()
	_, _ = fmt.Fprintf(hash, "%s%s%s", data.FirstName, data.LastName, data.DateOfBirth)
	return hash.Sum(nil)[:16]
}
