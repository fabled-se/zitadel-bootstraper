package zitadel

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewJWT(keyDataJson []byte, domain string) (string, error) {
	data := struct {
		KeyID  string `json:"keyId"`
		Key    string `json:"key"`
		UserId string `json:"userId"`
	}{}

	if err := json.Unmarshal(keyDataJson, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal keyDataJson: %w", err)
	}

	privatePem, _ := pem.Decode([]byte(data.Key))
	privateKey, err := x509.ParsePKCS1PrivateKey(privatePem.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to pars PKCS1 private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": data.UserId,
		"sub": data.UserId,
		"aud": "https://" + domain,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	})

	token.Header["kid"] = data.KeyID

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt with private key: %w", err)
	}

	return tokenString, nil
}
