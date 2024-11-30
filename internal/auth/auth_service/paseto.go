package auth_service

import (
	"auth-service/internal/models"
	"github.com/vk-rv/pvx"
	"golang.org/x/crypto/ed25519"
	"time"
)

type PasetoAuth interface {
	VerifyToken(token string) (*models.ServiceClaims, error)
	NewToken(data models.TokenData) (string, *models.ServiceClaims, error)
}

type pasetoAuth struct {
	privateKey *pvx.AsymSecretKey
	publicKey  *pvx.AsymPublicKey
}

func NewPaseto() PasetoAuth {

	pubKey, privKey, _ := ed25519.GenerateKey(nil)

	privateKey := pvx.NewAsymmetricSecretKey(privKey, pvx.Version4)
	publicKey := pvx.NewAsymmetricPublicKey(pubKey, pvx.Version4)

	return &pasetoAuth{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (pa *pasetoAuth) NewToken(data models.TokenData) (string, *models.ServiceClaims, error) {

	serviceClaims := &models.ServiceClaims{}

	iss := time.Now().UTC()
	exp := iss.Add(data.Duration)

	serviceClaims.IssuedAt = &iss
	serviceClaims.Expiration = &exp
	serviceClaims.Subject = data.Subject

	serviceClaims.AdditionalClaims = data.AdditionalClaims
	serviceClaims.Footer = data.Footer

	pv4 := pvx.NewPV4Public()

	authToken, err := pv4.Sign(
		pa.privateKey,
		serviceClaims,
		pvx.WithFooter(serviceClaims.Footer))
	if err != nil {
		return "", nil, err
	}

	claims, err := pa.VerifyToken(authToken)
	if err != nil {
		return "", nil, err
	}

	return authToken, claims, nil
}

func (pa *pasetoAuth) VerifyToken(token string) (*models.ServiceClaims, error) {
	pv4 := pvx.NewPV4Public()
	tk := pv4.Verify(token, pa.publicKey)

	f := models.Footer{}
	sc := models.ServiceClaims{
		Footer: f,
	}

	err := tk.Scan(&sc, &f)
	if err != nil {
		return &sc, err
	}

	return &sc, nil
}
