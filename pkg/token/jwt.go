package token

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
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

func (j *JWT) GenerateToken() (string, error) {
	// token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   j.Issuer,
		"exp":   time.Now().Add(time.Minute * 5).Unix(),
		"roles": []string{"POLYTUNE_GW_INTERNAL_ACCESS"},
	})

	return token.SignedString([]byte(j.Secret))
}

func (j *JWT) RestyMiddleware() resty.PreRequestHook {
	return func(c *resty.Client, req *http.Request) error {
		token, err := j.GenerateToken()
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}
