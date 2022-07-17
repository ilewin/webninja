package utils

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	App_Port        string `mapstructure:"APP_PORT"`
	App_Storage     string `mapstructure:"APP_STORAGE"`
	App_Files_Field string `mapstructure:"APP_FILES_FIELD"`

	Cors_Origin  string `mapstructure:"CORS_ORIGINS"`
	Cors_Header  string `mapstructure:"CORS_HEADER"`
	Cors_Methods string `mapstructure:"CORS_METHODS"`

	Email_To      string `mapstructure:"EMAIL_TO"`
	Email_From    string `mapstructure:"EMAIL_FROM"`
	Email_Subject string `mapstructure:"EMAIL_SUBJECT"`

	Smtp_Login  string `mapstructure:"SMTP_LOGIN"`
	Smtp_Pass   string `mapstructure:"SMTP_PASS"`
	Smtp_Host   string `mapstructure:"SMTP_HOST"`
	Smtp_Server string `mapstructure:"SMTP_SERVER"`

	Compression int `mapstructure:"COMPRESSION"`

	Slack_Auth_Token string `mapstructure:"SLACK_AUTH_TOKEN"`
	Slack_Channel_Id string `mapstructure:"SLACK_CHANNEL_ID"`
	Slack_App_Token  string `mapstructure:"SLACK_APP_TOKEN"`
}

func loadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

var app_config Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		var env_err error
		app_config, env_err = loadConfig(".")
		if env_err != nil {
			log.Fatal("cannot load config:", env_err)
			panic("Cant continue")
		}
	})
	return &app_config
}
