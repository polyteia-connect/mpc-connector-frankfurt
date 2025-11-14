package token

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Issuer string
	key    *ecdsa.PrivateKey
}

func NewJWT(keyFilePath string, issuer string) (*JWT, error) {
	// Read the key file
	keyData, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	// Parse the PEM-encoded ECDSA private key
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from key file")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		// Try parsing as PKCS8 format
		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		var ok bool
		key, ok = parsedKey.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not an ECDSA private key")
		}
	}

	return &JWT{
		Issuer: issuer,
		key:    key,
	}, nil
}

func (j *JWT) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss":   j.Issuer,
		"exp":   time.Now().Add(time.Minute * 5).Unix(),
		"roles": []string{"POLYTUNE_GW_INTERNAL_ACCESS"},
	})

	return token.SignedString(j.key)
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
