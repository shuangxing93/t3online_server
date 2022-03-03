package config

import "github.com/spf13/viper"

type Config struct {
	Database database
}

type database struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBTable       string `mapstructure:"DB_TABLE"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASS"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("dbconfig")
	viper.SetConfigType("toml")

	// viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
