package main

import (
	"github.com/spf13/viper"
	"github.com/vincentserpoul/gorestarter/pkg/storage"
)

// config is the app configuration
type config struct {
	MySQLDBConf *storage.MySQLDBConf
	HTTPPort    int
}

// newConfig will retrieve the current config
func newConfig() *config {

	viper.SetDefault("httpport", int(9002))
	viper.SetDefault("mysqldb", map[string]string{
		"protocol": "tcp",
		"host":     "127.0.0.1",
		"port":     "3306",
		"user":     "internal",
		"password": "dev",
		"dbname":   "dev",
	})

	return &config{
		MySQLDBConf: &storage.MySQLDBConf{
			Protocol: viper.GetStringMapString("mysqldb")["protocol"],
			Host:     viper.GetStringMapString("mysqldb")["host"],
			Port:     viper.GetStringMapString("mysqldb")["port"],
			User:     viper.GetStringMapString("mysqldb")["user"],
			Password: viper.GetStringMapString("mysqldb")["password"],
			DbName:   viper.GetStringMapString("mysqldb")["dbname"],
		},
		HTTPPort: viper.GetInt("httpport"),
	}
}
