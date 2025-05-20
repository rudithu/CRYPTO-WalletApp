package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	configMap map[string]string
	once      sync.Once
	configErr error
)

const (
	DB_USER  = "database.user"
	DB_PASS  = "database.password"
	DB_HOST  = "database.host"
	DB_PORT  = "database.port"
	DB_NAME  = "database.name"
	APP_PORT = "app.port"
)

func GetConfig() (map[string]string, error) {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")

		err := viper.ReadInConfig()
		if err != nil {
			log.Print("ERROR: Error reading config file:", err)
			configErr = fmt.Errorf("error reading config file: %s", err.Error())
			return
		}
		log.Println("config laoded")
		configMap = make(map[string]string)
		for _, key := range viper.AllKeys() {
			configMap[key] = viper.GetString(key)
		}
	})
	return configMap, configErr
}
