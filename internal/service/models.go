package service

import (
	"github.com/vk-rv/pvx"
	"time"
)

type AdditionalClaims struct {
	ID   string `json:"user_id"`
	Role string `json:"role"`
}

type Footer struct {
	MetaData string `json:"metadata"`
}

type ServiceClaims struct {
	pvx.RegisteredClaims
	AdditionalClaims
	Footer
}

type TokenData struct {
	Subject  string
	Duration time.Duration
	AdditionalClaims
	Footer
}
