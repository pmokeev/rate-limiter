package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"net"
	"strconv"
)

const (
	requestCount int32 = 100
	ipSubnet     int8  = 24
)

type Service struct {
	redisClient *redis.Client
}

func NewService() *Service {
	return &Service{redisClient: redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})}
}

func (s *Service) HaveAccess(context context.Context, IPv4 string) (bool, error) {
	_, IPNet, _ := net.ParseCIDR(IPv4 + "/" + strconv.Itoa(int(ipSubnet)))

	data, err := s.redisClient.Get(context, IPNet.String()).Result()
	if err == redis.Nil {
		return true, nil
	} else if err != nil {
		return false, err
	} else {
		var count int32
		if err := json.Unmarshal(bytes.NewBufferString(data).Bytes(), &count); err != nil {
			return false, nil
		}

		return count < requestCount, nil
	}
}

func (s *Service) GetData() {

}

func (s *Service) ClearRate() {

}
