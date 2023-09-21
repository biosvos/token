package main

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"os"
	"time"
)

type TokenCreator struct {
	private *rsa.PrivateKey
	method  jwt.SigningMethod
}

func NewTokenCreator(privateFile string) (*TokenCreator, error) {
	contents, err := os.ReadFile(privateFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(contents)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	method, err := newSigningMethod()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &TokenCreator{
		private: privateKey,
		method:  method,
	}, nil
}

func newSigningMethod() (jwt.SigningMethod, error) {
	const signingMethod = "RS256"
	alg := jwt.GetSigningMethod(signingMethod)
	if alg == nil {
		return nil, errors.New("not found signing method")
	}
	return alg, nil
}

func (t *TokenCreator) Create(data []byte, duration time.Duration) (string, error) {
	claims, err := newClaim(data, time.Now(), duration)
	if err != nil {
		return "", errors.WithStack(err)
	}
	token := jwt.NewWithClaims(t.method, claims)
	ret, err := token.SignedString(t.private)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return ret, nil
}

func newClaim(data []byte, now time.Time, duration time.Duration) (jwt.MapClaims, error) {
	var ret jwt.MapClaims
	err := json.Unmarshal(data, &ret)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ret["iat"] = now.UTC().Unix()
	ret["exp"] = now.UTC().Add(duration).Unix()
	return ret, nil
}
