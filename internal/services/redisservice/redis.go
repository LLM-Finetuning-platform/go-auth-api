package redisservice

import (
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Address string
	DB int
	Pwd string
	Protocol int

}


func (service *RedisService) Connect () (*redis.Client, error){
	newClient := redis.NewClient(&redis.Options{Addr: service.Address, DB: service.DB, Password: service.Pwd, Protocol: service.Protocol})
	defer service.disconnect(newClient)
	return newClient,nil
}

func (service *RedisService) disconnect (client *redis.Client) (error) {
	err:= client.Close()
	if err != nil {
		return err
	}
	return nil
}

