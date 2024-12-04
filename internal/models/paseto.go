package models

import (
	"github.com/vk-rv/pvx"
)

type AdditionalClaims struct {
	ID string
}

type Footer struct {
	MetaData string `json:"metadata"`
}

type ServiceClaims struct {
	pvx.RegisteredClaims
	AdditionalClaims
	Footer
}
