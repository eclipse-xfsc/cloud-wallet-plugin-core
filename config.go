package core

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

var libConfig Config

func init() {
	libConfig = getConfig()
}

func getConfig() Config {
	conf := Config{}
	err := envconfig.Process("PLUGIN", &conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

type Config struct {
	LogLevel string `default:"info"`
	IsDev    bool   `default:"true"`
	Name     string
	Tenant   string
	KeyCloak struct {
		Url       string
		Login     string
		Password  string
		RealmName string
		TokenTTL  time.Duration `default:"250ns"`
	}
	Policy struct {
		Url              string
		Repository       string
		Group            string
		ExpiresInSeconds int
	}
	Nats struct {
		Url        string
		QueueGroup string
	}
	Crypto struct {
		Namespace string
	}
	DIDComm struct {
		Url string
	}
}

func SetLibConfig(c Config) {
	libConfig = c
}
