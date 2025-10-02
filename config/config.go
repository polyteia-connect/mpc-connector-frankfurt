package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port              int      `env:"PORT" env-default:"3000"`
	Host              string   `env:"HOST" env-default:"0.0.0.0"`
	ExternalClientURL string   `env:"EXTERNAL_CLIENT_URL"`
	Debug             bool     `env:"DEBUG" env-default:"false"`
	MPCBaseURL        string   `env:"MPC_BASE_URL" env-default:"http://localhost:8000"`
	MPCLeaderID       int      `env:"MPC_LEADER_ID" env-default:"0"`
	MPCPartyID        int      `env:"MPC_PARTY_ID" env-default:"0"`
	MPCParticipants   []string `env:"MPC_PARTICIPANTS" env-required:"true"`
}

func Auto() Config {
	var config Config

	if err := cleanenv.ReadConfig(".env", &config); err != nil {
		if err := cleanenv.ReadEnv(&config); err != nil {
			panic(fmt.Errorf("failed to read config: %w", err))
		}
	}

	return config
}
