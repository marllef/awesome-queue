package configs

type configs struct {
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort int `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDatabase int `mapstructure:"REDIS_DATABASE"`
	RedisConsumerID string `mapstructure:"REDIS_CONSUMER_ID"`
}


