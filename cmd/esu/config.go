package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port               int      `env:"PORT" env-default:"3000"`
	Host               string   `env:"HOST" env-default:"0.0.0.0"`
	VaccinationBaseURL string   `env:"VACCINATION_BASE_URL" env-required:"true"`
	Debug              bool     `env:"DEBUG" env-default:"false"`
	JWTKeyFile         string   `env:"JWT_KEY_FILE" env-required:"true"`
	JWTIssuer          string   `env:"JWT_ISSUER" env-required:"true"`
	MPCBaseURL         string   `env:"MPC_BASE_URL" env-required:"true"`
	MPCLeaderID        int      `env:"MPC_LEADER_ID" env-default:"0"`
	MPCPartyID         int      `env:"MPC_PARTY_ID" env-default:"1"`
	MPCParticipants    []string `env:"MPC_PARTICIPANTS" env-required:"true"`
}

func AutoConfig() Config {
	var config Config

	if err := cleanenv.ReadConfig(".env", &config); err != nil {
		if err := cleanenv.ReadEnv(&config); err != nil {
			panic(fmt.Errorf("failed to read config: %w", err))
		}
	}

	return config
}
