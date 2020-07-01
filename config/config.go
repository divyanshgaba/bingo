package config

import (
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	env   = flag.String("env", "prod", "deployment environment for config. Default `prod`")
	mongo *MongoConfiguration
)

func init() {
	flag.Parse()
	viper.AddConfigPath("$GOPATH/src/github.com/divyanshgaba/bingo")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetConfigName("config-" + *env)
	if err := viper.ReadInConfig(); err != nil {
		log.WithField("err", err).Fatal("couldn't read base config")
	}
	if err := viper.Sub("mongo").Unmarshal(&mongo); err != nil {
		log.WithField("err", err).Fatal("couldn't marshalling mongo config")
	}
}

// Mongo returns instance of MongoConfiguration
func Mongo() MongoConfiguration {
	return *mongo
}

// Env returns environment that is used for config
func Env() string {
	return *env
}
