package configs

import (
	"flag"

	"github.com/marllef/awesome-queue/pkg/utils/path"
	"github.com/spf13/viper"
)

var envPath = flag.String("env-path", path.Root, "Defines the path of the env file")
var envName = flag.String("env-name", "local", "Defines the name of the env file")
var envType = flag.String("env-type", "env", "Defines the type of the env file")
var Env *configs

// Load environment variables. 
func LoadEnvs() (err error) {
	if Env == nil {
		Env, err = loadEnvFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func loadEnvFile() (config *configs, err error) {
	viper.AddConfigPath(*envPath)
	viper.SetConfigName(*envName)
	viper.SetConfigType(*envType)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
