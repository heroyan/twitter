package config

import "github.com/spf13/viper"

func init() {
	viper.SetConfigFile("./config.json")
	viper.SetConfigType("json")
}

func GetAddr() string {
	return viper.GetString("addr")
}

func GetPasswd() string {
	return viper.GetString("passwd")
}

func GetDB() int {
	return viper.GetInt("db")
}

func GetSessionKey() string {
	return viper.GetString("session_key")
}

func GetSessionAge() int {
	return viper.GetInt("session_expire")
}
