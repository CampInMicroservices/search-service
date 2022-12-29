package config

import (
	"log"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	Environment   string `mapstructure:"ENVIRONMENT"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	ConsulAddress string `mapstructure:"CONSUL_ADDRESS"`
	LogitAddress  string `mapstructure:"LOGIT_ADDRESS"`
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	MigrationURL  string `mapstructure:"MIGRATION_URL"`
	GinMode       string `mapstructure:"GIN_MODE"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func WatchConsulConfig(key string, consulAddress string, callback func(source string)) {

	consul_config := api.DefaultConfig()
	consul_config.Address = consulAddress
	client, err := api.NewClient(consul_config)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	kv := client.KV()

	_, meta, err := kv.Get(key, nil)
	if err != nil {
		return
	}

	options := api.QueryOptions{
		RequireConsistent: true,
	}

	var pair *api.KVPair
	for {
		options.WaitIndex = meta.LastIndex
		pair, meta, err = kv.Get(key, &options)

		if err != nil {
			log.Panic("Key not available: ", err)
			return
		}

		if pair != nil {
			callback(string(pair.Value))
		}
	}
}
