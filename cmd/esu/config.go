package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port               int      `env:"PORT" env-default:"3000"`
	Host               string   `env:"HOST" env-default:"0.0.0.0"`
	VaccinationBaseURL string   `env:"VACCINATION_BASE_URL"`
	Debug              bool     `env:"DEBUG" env-default:"false"`
	JWTTokenSecret     string   `env:"JWT_TOKEN_SECRET"`
	JWTIssuer          string   `env:"JWT_ISSUER"`
	MPCBaseURL         string   `env:"MPC_BASE_URL"`
	MPCLeaderID        int      `env:"MPC_LEADER_ID"`
	MPCPartyID         int      `env:"MPC_PARTY_ID"`
	MPCParticipants    []string `env:"MPC_PARTICIPANTS"`
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
