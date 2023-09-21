package main

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"os"
)

type TokeParser struct {
	public *rsa.PublicKey
}

func NewTokeParser(publicFile string) (*TokeParser, error) {
	crt, err := os.ReadFile(publicFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	public, err := jwt.ParseRSAPublicKeyFromPEM(crt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &TokeParser{
		public: public,
	}, nil
}

func (t *TokeParser) Parse(token string) (map[string]any, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return t.public, nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	claims := parsed.Claims.(jwt.MapClaims)
	return claims, nil
}
