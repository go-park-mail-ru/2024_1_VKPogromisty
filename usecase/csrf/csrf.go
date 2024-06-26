package csrf

import (
	"os"
	"socio/errors"
	customtime "socio/pkg/time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	cryptAlg = "HS256"
)

type CSRFService struct {
	secret []byte
	TP     customtime.TimeProvider
}

func NewCSRFService(tp customtime.TimeProvider) (service *CSRFService) {
	return &CSRFService{
		secret: []byte(os.Getenv("CSRF_SECRET")),
		TP:     tp,
	}
}

type JwtCsrfClaims struct {
	SessionID string `json:"sid"`
	UserID    uint   `json:"uid"`
	jwt.StandardClaims
}

func (c *CSRFService) Create(sessionID string, userID uint, tokenExpTime int64) (token string, err error) {
	data := JwtCsrfClaims{
		SessionID: sessionID,
		UserID:    userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpTime,
			IssuedAt:  c.TP.Now().Unix(),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString(c.secret)
	if err != nil {
		return
	}

	return
}

func (c *CSRFService) parseSecretGetter(token *jwt.Token) (res interface{}, err error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)

	if !ok || method.Alg() != cryptAlg {
		err = errors.ErrInvalidJWT
		return
	}

	res = c.secret

	return
}

func (c *CSRFService) Check(sessionID string, userID uint, inputToken string) (err error) {
	payload := &JwtCsrfClaims{}

	_, err = jwt.ParseWithClaims(inputToken, payload, c.parseSecretGetter)
	if err != nil {
		err = errors.ErrInvalidJWT
		return
	}

	if payload.Valid() != nil {
		err = errors.ErrInvalidJWT
		return
	}

	if payload.SessionID != sessionID || payload.UserID != userID {
		err = errors.ErrInvalidJWT
		return
	}

	return
}
