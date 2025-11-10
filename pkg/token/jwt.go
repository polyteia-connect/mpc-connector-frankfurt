package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret string
	Issuer string
}

func NewJWT(secret string, issuer string) *JWT {
	return &JWT{
		Secret: secret,
		Issuer: issuer,
	}
}

func (j *JWT) GenerateToken(subject string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub":   subject,
		"iss":   j.Issuer,
		"exp":   time.Now().Add(time.Minute * 5).Unix(),
		"roles": []string{"POLYTUNE_GW_INTERNAL_ACCESS"},
	})

	return token.SignedString([]byte(j.Secret))
}
