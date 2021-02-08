package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mariusmagureanu/burlan/src/pkg/errors"

	"time"
)

// JwtWrapper wraps the signing key and the issuer
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token
type JwtClaim struct {
	jwt.StandardClaims
	Email     string
	Name      string
	ClientUID string
}

// GenerateToken generates a jwt token
func (j *JwtWrapper) GenerateToken(uid, name, email string) (string, error) {
	claims := &JwtClaim{
		Name:      name,
		ClientUID: uid,
		Email:     email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(j.SecretKey))
}

//ValidateToken validates the jwt token
func (j *JwtWrapper) ValidateToken(signedToken string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		return nil, errors.ErrCannotParseClaims
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.ErrJWTIsExpired
	}

	return claims, nil
}
