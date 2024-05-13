package util

import (
	"github.com/spf13/viper"
)


// Config stores the configuration for the application.
type Config struct {
	DBDriver      string `mapstructure:"POSTGRES_DRIVER"`
	DBSource      string `mapstructure:"POSTGRES_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	MongoDBName   string `mapstructure:"MONGO_DB_NAME"`
	MongoRootUser string `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	MongoRootPass string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`
	MongoUser     string `mapstructure:"MONGO_USER"`
	MongoPass     string `mapstructure:"MONGO_PASSWORD"`
	MongoSource   string `mapstructure:"MONGO_SOURCE"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	viper.AddConfigPath(path)
	viper.SetConfigName("")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, err
}
