package redisservice

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Address string
	DB int
	Pwd string
	Protocol int
	Poolsize int
	MinIdleConnections int
	PoolTimeOutSecs int
}


func (service *RedisService) Connect () (*redis.Client, error){
	newClient := redis.NewClient(&redis.Options{Addr: service.Address, DB: service.DB, Password: service.Pwd, Protocol: service.Protocol, 
	PoolSize: service.Poolsize,
	MinIdleConns: service.MinIdleConnections,
	PoolTimeout: time.Second*time.Duration(service.PoolTimeOutSecs),})
	return newClient,nil
}

func (service *RedisService) Disconnect (client *redis.Client) (error) {
	err:= client.Close()
	if err != nil {
		return err
	}
	return nil
}

