package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

var Conf Config

type Config struct {
	Port                   uint     `mapstructure:"PORT"`
	Database               Database `mapstructure:",squash"`
	CronSpecsToFetchVideos string   `mapstructure:"CRON_TO_FETCH_VIDEOS"`
	VideoKeyword           string   `mapstructure:"KEYWORD_TO_FETCH_VIDEOS"`
	GoogleAPIKeys          []string
	ActiveGoogleAPIKey     string
	Ampq                   Amqp `mapstructure:",squash"`
}

type Amqp struct {
	Url       string `mapstructure:"AMQP_URL"`
	QueueName string `mapstructure:"AMQP_QUEUE_NAME"`
}

type Database struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     uint   `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	SSLMode  string `mapstructure:"DB_SSL_MODE"`
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
	// convert the google api keys to a slice of strings
	Conf.GoogleAPIKeys = strings.Split(googleAPIKeys.(string), ",")
	// by default set first key as active api key
	Conf.ActiveGoogleAPIKey = Conf.GoogleAPIKeys[0]
	return
}
