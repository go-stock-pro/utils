package token

import (
	"context"
	"crypto/rsa"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/rohanraj7316/utils/constants"
	"github.com/rohanraj7316/utils/env"
	"github.com/rohanraj7316/utils/redisclient"
)

type Handler struct {
	authPublicKey  *rsa.PublicKey
	authPrivateKey *rsa.PrivateKey
	location       *time.Location
	rh             redisclient.Handler
}

func New() (h *Handler, err error) {

	rEnv := []string{
		ENV_AUTH_PUBLIC_KEY,
		ENV_AUTH_PRIVATE_KEY,
	}
	uEnv := []string{}
	eVar := env.EnvData(rEnv, uEnv)

	pvtKey, err := LoadPrivateKey(derToPemString(eVar[ENV_AUTH_PRIVATE_KEY], "PRIVATE KEY"), "")
	if err != nil {
		return h, err
	}

	pubKey, err := LoadPublicKey(derToPemString(eVar[ENV_AUTH_PUBLIC_KEY], "PUBLIC KEY"))
	if err != nil {
		return h, err
	}

	loc, err := time.LoadLocation(constants.LOCATION)
	if err != nil {
		return h, err
	}

	rh, err := redisclient.New()
	if err != nil {
		return h, err
	}

	h = &Handler{
		authPublicKey:  pubKey,
		authPrivateKey: pvtKey,
		location:       loc,
		rh:             rh,
	}

	return h, nil
}

func (h *Handler) GenerateAuthToken(ctx context.Context, userId string) (res GenerateAuthTokenResponse, err error) {
	expireAt := h.eod().Format(constants.TIME_FORMAT)
	claims := jwt.MapClaims{
		"id":  userId,
		"eAt": expireAt,
		"cAt": time.Now().In(h.location).Format(constants.TIME_FORMAT),
	}

	j := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	token, err := j.SignedString(h.authPrivateKey)
	if err != nil {
		return res, err
	}

	res = GenerateAuthTokenResponse{
		AuthToken: token,
		ExpireAt:  expireAt,
	}

	return res, nil
}

func (h *Handler) ValidateAuthToken(ctx context.Context, token string) error {

	vt, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return h.authPublicKey, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, ok := vt.Claims.(jwt.MapClaims)
	if !ok || !vt.Valid {
		return errors.Errorf("failed to parse token claims")
	}

	// claims, ok := vt.Claims.(jwt.MapClaims)
	// if !ok || !vt.Valid {
	// 	return errors.Errorf("failed to parse token claims")
	// }

	// if id, ok := claims["id"].(string); !ok {
	// 	return errors.Errorf("failed to parse user id")
	// }

	// now := time.Now()
	// if ex, ok := claims["eAt"].(string); !ok {
	// 	return errors.Errorf("failed to parse expire at to string")
	// } else {
	// 	eAt, err := time.Parse(constants.TIME_FORMAT, ex)
	// 	if err != nil {
	// 		return errors.Errorf("failed to parse expire at to time")
	// 	}

	// 	if eAt.Before(now) {
	// 		return errors.Errorf("expired auth token")
	// 	}
	// }

	return nil
}
