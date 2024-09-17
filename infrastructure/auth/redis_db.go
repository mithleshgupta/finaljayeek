package auth

import (
	"github.com/go-redis/redis/v9"
)

// RedisService struct contains the redis client, and auth service
type RedisService struct {
	RedisClient *redis.Client
	AuthService AuthServiceInterface
}

// NewRedisService creates a new instance of RedisService with the provided host, port, and password
func NewRedisService(host, port, password string) (*RedisService, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return &RedisService{
		RedisClient: redisClient,
		AuthService: NewAuthService(redisClient),
	}, nil
}
