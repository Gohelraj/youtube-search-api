package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

var Conf Config

type Config struct {
	Port                   string `mapstructure:"PORT"`
	DatabaseURL            string `mapstructure:"DATABASE_URL"`
	CronSpecsToFetchVideos string `mapstructure:"CRON_TO_FETCH_VIDEOS"`
	VideoKeyword           string `mapstructure:"KEYWORD_TO_FETCH_VIDEOS"`
	GoogleAPIKeys          []string
	Ampq                   Amqp `mapstructure:",squash"`
}

type Amqp struct {
	Url       string `mapstructure:"AMQP_URL"`
	QueueName string `mapstructure:"AMQP_QUEUE_NAME"`
}

// LoadConfig loads the config variables and returns a config struct
func LoadConfig() (err error) {
	// set config file's path, name and type
	viper.SetConfigFile(".env")
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("no config file found. going on...")
		} else {
			log.Printf("error loading config file")
		}
		return
	}
	// load environment variables in the config struct
	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Printf("error while unmarshaling the config into a Struct: %v\n", err)
		return
	}
	googleAPIKeys := viper.Get("GOOGLE_API_KEYS")
	Conf.GoogleAPIKeys = strings.Split(googleAPIKeys.(string), ",")
	return
}
