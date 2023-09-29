package redis

import (
	redis "github.com/redis/go-redis/v9"
	"github.com/zilliangroup/builder-backend/src/utils/config"
	"go.uber.org/zap"
)

const RETRY_TIMES = 6

type RedisConfig struct {
	Addr     string `env:"ZWEB_REDIS_ADDR" envDefault:"localhost"`
	Port     string `env:"ZWEB_REDIS_PORT" envDefault:"6379"`
	Password string `env:"ZWEB_REDIS_PASSWORD" envDefault:"zweb2022"`
	Database int    `env:"ZWEB_REDIS_DATABASE" envDefault:"0"`
}

func NewRedisConnectionByGlobalConfig(config *config.Config, logger *zap.SugaredLogger) (*redis.Client, error) {
	redisConfig := &RedisConfig{
		Addr:     config.GetRedisAddr(),
		Port:     config.GetRedisPort(),
		Password: config.GetRedisPassword(),
		Database: config.GetRedisDatabase(),
	}
	return NewRedisConnection(redisConfig, logger)
}

func NewRedisConnection(config *RedisConfig, logger *zap.SugaredLogger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr + ":" + config.Port,
		Password: config.Password,
		DB:       config.Database,
	})

	logger.Infow("connected with redis", "redis", config)

	return rdb, nil
}
