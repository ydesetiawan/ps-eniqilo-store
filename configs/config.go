package configs

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type MainConfig struct {
	Host         string `env:"DB_HOST"`
	Port         string `env:"DB_PORT"`
	Username     string `env:"DB_USERNAME"`
	Password     string `env:"DB_PASSWORD"`
	DatabaseName string `env:"DB_NAME"`
	Params       string `env:"DB_PARAMS"`
}

func (dc *MainConfig) GetDsnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", dc.Username, dc.Password, dc.Host, dc.Port, dc.DatabaseName, dc.Params)
}

func (mc *MainConfig) ReadConfig() error {
	v := reflect.ValueOf(mc).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("env")
		if tag != "" {
			envValue := os.Getenv(tag)
			if envValue != "" {
				field.SetString(envValue)
			}
		}
	}
	return nil
}

func Init() *MainConfig {
	godotenv.Load(".env")

	config := &MainConfig{}
	config.ReadConfig()
	return config
}
