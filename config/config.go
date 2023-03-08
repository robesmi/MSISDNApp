package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/robesmi/MSISDNApp/model/errs"
	"github.com/spf13/viper"
)

type Config struct {
	MySqlDriver           	string	`mapstructure:"MYSQL_DRIVER"`
	MySqlSource           	string	`mapstructure:"MYSQL_SOURCE"`
	Port                  	string	`mapstructure:"PORT"`
	AccessTokenPrivateKey 	string	`mapstructure:"AccessTokenPrivateKey"`
	AccessTokenPublicKey  	string	`mapstructure:"AccessTokenPublicKey"`
	RefreshTokenPrivateKey 	string	`mapstructure:"RefreshTokenPrivateKey"`
	RefreshTokenPublicKey  	string	`mapstructure:"RefreshTokenPublicKey"`
	GoogleClientID			string	`mapstructure:"GoogleClientID"`
	GoogleClientSecret		string	`mapstructure:"GoogleClientSecret"`
	GithubClientID			string	`mapstructure:"GithubClientID"`
	GithubClientSecret		string	`mapstructure:"GithubClientSecret"`
}

func LoadConfig() (Config, *errs.AppError){


	_, b, _, _ := runtime.Caller(0)
	basepath   := filepath.Dir(b)
	
	var config Config
	viper.AddConfigPath(basepath)
	viper.SetConfigType("env")
	viper.SetConfigName("config")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed reading config " + err.Error())
		return config, errs.UnexpectedError(err.Error())
	}
	err = viper.Unmarshal(&config)
	if err != nil{
		log.Fatalf("Failed unmarshaling config")
		return config, errs.UnexpectedError(err.Error())
	}
	return config, nil

}