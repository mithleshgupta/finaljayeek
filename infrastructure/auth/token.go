package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Token struct implements the TokenInterface
type Token struct{}

// TokenInterface defines the methods for the token struct
type TokenInterface interface {
	CreateToken(userID uint64) (*TokenDetails, error)
	ExtractTokenMetadata(*http.Request) (*AccessDetails, error)
}

// NewToken creates a new instance of the Token struct
func NewToken() *Token {
	return &Token{}
}

// CreateToken creates a new token for the provided user ID
func (t *Token) CreateToken(userID uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AccessTokenExpiresAt = time.Now().Add(time.Minute * 15).Unix()
	td.AccessTokenUUID = uuid.New().String()

	td.RefreshTokenExpiresAt = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshTokenUUID = td.AccessTokenUUID + "++" + strconv.Itoa(int(userID))

	var err error
	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_token_uuid"] = td.AccessTokenUUID
	atClaims["user_id"] = userID
	atClaims["expires_at"] = td.AccessTokenExpiresAt
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	// Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_token_uuid"] = td.RefreshTokenUUID
	rtClaims["user_id"] = userID
	rtClaims["expires_at"] = td.RefreshTokenExpiresAt
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// TokenValid checks if the token in the request is valid
func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return err
	}
	return nil
}

// VerifyToken parses and verifies the token in the request
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ExtractToken gets the token from the request header
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// ExtractTokenMetadata extracts the token metadata from the request
func (t *Token) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessTokenUUID, ok := claims["access_token_uuid"].(string)
		if !ok {
			return nil, fmt.Errorf("access_token_uuid not found in claims")
		}
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("user_id not found in claims")
		}
		return &AccessDetails{
			AccessTokenUUID: accessTokenUUID,
			UserID:          uint64(userID),
		}, nil
	}
	return nil, fmt.Errorf("invalid token")
}
