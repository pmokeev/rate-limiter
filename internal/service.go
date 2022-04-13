package internal

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	requestCount int64         = 100
	ipSubnet     int64         = 24
	timeToLive   time.Duration = 120 * time.Second
)

type Service struct {
	redisClient *redis.Client
}

func NewService() *Service {
	return &Service{redisClient: redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})}
}

func (s *Service) CheckAccess(context context.Context, IPv4 string) (*middlewareAccess, error) {
	_, IPNet, _ := net.ParseCIDR(IPv4 + "/" + strconv.Itoa(int(ipSubnet)))

	count, err := s.redisClient.Incr(context, IPNet.String()).Result()
	if err == redis.Nil {
		if err := s.redisClient.Set(context, IPNet.String(), 1, timeToLive).Err(); err != nil {
			return NewMiddlewareAccess(0, requestCount, 0, false), err
		}

		return NewMiddlewareAccess(requestCount-1, requestCount, 0, true), nil
	} else if err != nil {
		return NewMiddlewareAccess(0, requestCount, 0, false), err
	} else {
		if err != nil {
			return NewMiddlewareAccess(0, requestCount, 0, false), err
		}

		ttl, err := s.redisClient.TTL(context, IPNet.String()).Result()
		if err != nil {
			return NewMiddlewareAccess(0, requestCount, 0, false), err
		}

		return NewMiddlewareAccess(requestCount-count, requestCount, ttl, count < requestCount), nil
	}
}

func (s *Service) GetData() (string, error) {
	response, err := http.Get("https://api.chucknorris.io/jokes/random")
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var data chuckJoke
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.Value, nil
}

func (s *Service) ClearRate(context context.Context, IPv4 string) error {
	_, IPNet, _ := net.ParseCIDR(IPv4 + "/" + strconv.Itoa(int(ipSubnet)))
	_, err := s.redisClient.Del(context, IPNet.String()).Result()

	return err
}
