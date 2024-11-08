package auth

import (
	"github.com/pkg/errors"
	"github.com/vk-rv/pvx"
	"time"
)

type PasetoAuth struct {
	privateKey    *pvx.AsymSecretKey
	publicKey     *pvx.AsymPublicKey
	asymmetricKey []byte
}

const keySize = 32

var (
	ErrInvalidSize = errors.Errorf("bad key size: it must be %d bytes", keySize)
)

func NewPaseto(key []byte) (*PasetoAuth, error) {

	if len(key) != keySize {
		return nil, ErrInvalidSize
	}

	privateKey := pvx.NewAsymmetricSecretKey(key, pvx.Version4)
	publicKey := pvx.NewAsymmetricPublicKey(key, pvx.Version4)

	return &PasetoAuth{
		asymmetricKey: key,
		privateKey:    privateKey,
		publicKey:     publicKey,
	}, nil
}

func (pa *PasetoAuth) NewToken(data TokenData) (string, error) {

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

func (pa *PasetoAuth) VerifyToken(token string) (*ServiceClaims, error) {
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
