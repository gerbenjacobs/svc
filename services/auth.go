package services

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
)

// Auth is the service tasked with creating and validating jwt tokens
type Auth struct {
	secretToken []byte
}

func NewAuth(secretToken []byte) *Auth {
	return &Auth{
		secretToken: secretToken,
	}
}

// UserClaims is a custom struct used as jwt payload
type UserClaims struct {
	UserID string
	*jwt.RegisteredClaims
}

func (u UserClaims) Valid() error {
	if u.UserID == "" {
		return errors.New("user id missing from claims")
	}
	return nil
}

func (a *Auth) tokenFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return a.secretToken, nil
}

func (a *Auth) Create(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserID: userID,
		RegisteredClaims: &jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
		},
	})

	return token.SignedString(a.secretToken)
}

func (a *Auth) ReadFromRequest(r *http.Request) (*UserClaims, error) {
	ext := request.MultiExtractor([]request.Extractor{
		request.OAuth2Extractor,
	})
	token, err := request.ParseFromRequest(r, ext, a.tokenFunc, request.WithClaims(&UserClaims{}))
	if err != nil {
		return nil, fmt.Errorf("failed to parse authentication request: %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid user claims")
	}

	if !claims.VerifyNotBefore(time.Now().UTC(), false) {
		return nil, fmt.Errorf("token not valid before: %v", claims.NotBefore)
	}

	return claims, nil
}
