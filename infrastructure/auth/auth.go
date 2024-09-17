package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
)

var (
	ctx = context.Background()
)

// AuthServiceInterface defines the methods that an authentication service should implement.
type AuthServiceInterface interface {
	CreateAuth(uint64, *TokenDetails) error
	FetchAuth(string) (uint64, error)
	DeleteRefreshToken(string) error
	DeleteTokens(*AccessDetails) error
}

// AuthService represents the authentication service implementation
type AuthService struct {
	redisClient *redis.Client
}

var _ AuthServiceInterface = &AuthService{}

// NewAuthService creates and returns a new instance of AuthService
func NewAuthService(redisClient *redis.Client) *AuthService {
	return &AuthService{redisClient: redisClient}
}

// AccessDetails represents the access details of a user.
type AccessDetails struct {
	AccessTokenUUID string
	UserID          uint64
}

// TokenDetails represents the token details of a user.
type TokenDetails struct {
	AccessToken           string
	RefreshToken          string
	AccessTokenUUID       string
	RefreshTokenUUID      string
	AccessTokenExpiresAt  int64
	RefreshTokenExpiresAt int64
}

type PasswordResetRequest struct {
	VerificationCode     string `json:"verification_code" validate:"required"`
	Phone                string `json:"phone" validate:"required,e164"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type PhoneVerificationRequest struct {
	Code  string `json:"code" validate:"required"`
	Phone string `json:"phone" validate:"required,e164"`
}

// CreateAuth saves token metadata to Redis
func (s *AuthService) CreateAuth(userID uint64, td *TokenDetails) error {
	accessTokenExpiresAt := time.Unix(td.AccessTokenExpiresAt, 0)
	refreshTokenExpiresAt := time.Unix(td.RefreshTokenExpiresAt, 0)
	now := time.Now()

	// set the token metadata in redis
	accessToken, err := s.redisClient.Set(ctx, td.AccessTokenUUID, strconv.Itoa(int(userID)), accessTokenExpiresAt.Sub(now)).Result()
	if err != nil {
		return err
	}
	refreshToken, err := s.redisClient.Set(ctx, td.RefreshTokenUUID, strconv.Itoa(int(userID)), refreshTokenExpiresAt.Sub(now)).Result()
	if err != nil {
		return err
	}
	// check if the metadata was saved correctly
	if accessToken != "OK" || refreshToken != "OK" {
		return errors.New("no record inserted")
	}
	return nil
}

// FetchAuth checks the metadata saved
func (s *AuthService) FetchAuth(tokenUUID string) (uint64, error) {
	// get the userID from the redis
	userID, err := s.redisClient.Get(ctx, tokenUUID).Result()
	if err != nil {
		return 0, err
	}
	userIDInt, _ := strconv.ParseUint(userID, 10, 64)
	return userIDInt, nil
}

// DeleteTokens deletes the access and refresh tokens of a user
func (s *AuthService) DeleteTokens(authD *AccessDetails) error {
	//get the refresh UUID
	refreshAccessTokenUUID := fmt.Sprintf("%s++%d", authD.AccessTokenUUID, authD.UserID)

	// delete access token
	deletedAccessToken, err := s.redisClient.Del(ctx, authD.AccessTokenUUID).Result()
	if err != nil {
		return err
	}
	// delete refresh token
	deletedRefreshToken, err := s.redisClient.Del(ctx, refreshAccessTokenUUID).Result()
	if err != nil {
		return err
	}
	// when the record is deleted, the return value is 1
	if deletedAccessToken != 1 || deletedRefreshToken != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

// DeleteRefresh deletes a refresh token
func (s *AuthService) DeleteRefreshToken(refreshTokenUUID string) error {
	// delete refresh token
	deletedRefreshToken, err := s.redisClient.Del(ctx, refreshTokenUUID).Result()
	if err != nil || deletedRefreshToken == 0 {
		return err
	}
	return nil
}
