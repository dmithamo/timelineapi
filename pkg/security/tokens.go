package security

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dmithamo/timelineapi/pkg/utils"
)

var signingKey = []byte(os.Getenv("SECRET"))

// CustomClaims standardizes the shape of custom claims
type CustomClaims struct {
	UID interface{} `json:"uuid"`
	jwt.StandardClaims
}

// GenerateToken generates a token with <any> claims embedded
func GenerateToken(claims interface{}) (*string, error) {
	tokenExpiryDate := time.Now().Add(10 * time.Minute).Unix()

	structuredClaims := CustomClaims{
		claims,
		jwt.StandardClaims{
			ExpiresAt: tokenExpiryDate,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, structuredClaims)

	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}
	return &signedToken, nil
}

// ValidateToken parses a token to validate it
func ValidateToken(signedToken string) (interface{}, error) {

	token, err := jwt.ParseWithClaims(signedToken, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	// getClaimsHelper extracts claims from token
	getClaimsHelper := func() (*CustomClaims, error) {
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			return nil, fmt.Errorf("invalid auth token")
		}

		return claims, nil
	}

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, jwt.ErrSignatureInvalid
		}

		errType := err.(*jwt.ValidationError)
		if errType.Errors == jwt.ValidationErrorExpired {
			// token exists, isValid(ish), but need refreshing
			claims, err := getClaimsHelper()
			if err != nil {
				return nil, err
			}
			return claims.UID, fmt.Errorf(utils.TOKEN_EXPIRED_ERR)
		}

		return nil, err
	}

	claims, err := getClaimsHelper()
	if err != nil {
		return nil, err
	}

	return claims.UID, err
}

// DecodeToken reads the cookie vlaue from a request and parses it
func DecodeToken(r *http.Request) (interface{}, error) {

	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}

	claims, err := ValidateToken(cookie.Value)
	if err != nil {
		if err.Error() != utils.TOKEN_EXPIRED_ERR {
			return nil, err
		}

		// refresh token, just to read its innards
		rt, err := GenerateToken(claims)

		if err != nil {
			return nil, err
		}

		claims, err = ValidateToken(*rt)
		if err != nil {
			return nil, err
		}
	}

	return claims, nil
}
