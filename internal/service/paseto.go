package service

import (
	"github.com/vk-rv/pvx"
	"golang.org/x/crypto/ed25519"
	"time"
)

type PasetoAuth interface {
	VerifyToken(token string) (*ServiceClaims, error)
	NewToken(data TokenData) (string, error)
}

type pasetoAuth struct {
	privateKey    *pvx.AsymSecretKey
	publicKey     *pvx.AsymPublicKey
	asymmetricKey []byte
}

func NewPaseto(key []byte) (PasetoAuth, error) {

	pubKey, privKey, _ := ed25519.GenerateKey(nil)

	privateKey := pvx.NewAsymmetricSecretKey(privKey, pvx.Version4)
	publicKey := pvx.NewAsymmetricPublicKey(pubKey, pvx.Version4)

	return &pasetoAuth{
		asymmetricKey: key,
		privateKey:    privateKey,
		publicKey:     publicKey,
	}, nil
}

func (pa *pasetoAuth) NewToken(data TokenData) (string, error) {

	serviceClaims := &ServiceClaims{}

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
		return "", err
	}

	return authToken, nil
}

func (pa *pasetoAuth) VerifyToken(token string) (*ServiceClaims, error) {
	pv4 := pvx.NewPV4Public()
	tk := pv4.Verify(token, pa.publicKey)

	f := Footer{}
	sc := ServiceClaims{
		Footer: f,
	}

	err := tk.Scan(&sc, &f)
	if err != nil {
		return &sc, err
	}

	return &sc, nil
}
