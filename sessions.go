package sessions

import (
	"time"

	"github.com/FlyCynomys/tools/log"
	redis "github.com/go-redis/redis"
)

type SessionCondig struct {
	CollectionTime int64
	SessionName    string
	Address        string
}

type SessionContext struct {
	Store       Store
	Cfg         *SessionCondig
	RedisClient *redis.Client
}

func (s *SessionContext) GetSessionContect() {

}
func (s *SessionContext) Init() {
	s.RedisClient = redis.NewClient(&redis.Options{
		Addr:         s.Cfg.Address,
		Password:     "",
		DB:           0,
		MaxRetries:   3,
		DialTimeout:  30,
		ReadTimeout:  10,
		WriteTimeout: 10,
		PoolSize:     10,
	})
	go func() {
		ti := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-ti.C:
				_, err := s.RedisClient.Ping().Result()
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()
}

func (s *SessionContext) GC() {

}

func (s *SessionContext) Read(sessionid string) (string, error) {
	value, err := s.RedisClient.Get(sessionid).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (s *SessionContext) SetSessionValue(sessionid string, value string, expire int64) error {
	err := s.RedisClient.Set(sessionid, value, time.Duration(expire)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionContext) DeleteSession(sessionid string) error {
	err := s.RedisClient.Del(sessionid).Err()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	return nil
}

func (s *SessionContext) ExpireSession(sessionid string, expire int64) error {
	err := s.RedisClient.Expire(sessionid, time.Duration(expire)).Err()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	return nil
}

type SesssionManage struct {
	SessionMap map[string]*SessionContext
}
