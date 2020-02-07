package utils

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var conf *Config

// Source by default
const defaultSource string = "https://data.gov.ru/opendata/7704786030-taxiparking/data-20140828T0000.json?encoding=UTF-8"
const defaultPort string = "8080"
const defaultRedis string = "redis:6379"

// defaultLimit limit for result
const defaultLimit int = 10

// Config configuration application
type Config struct {
	Port     string
	Redis    string
	Password string
	Db       int
	Source   string
}

// Configure application
func init() {
	conf = &Config{}
	flag.StringVar(&conf.Port, "port", "", "web server expose port")
	flag.StringVar(&conf.Redis, "redis", "", "host:port redis server")
	flag.StringVar(&conf.Password, "password", "", "password to redis")
	flag.StringVar(&conf.Source, "source", "", "url or path to source data")
	flag.IntVar(&conf.Db, "db", 0, "redis db index")
}

func parseArgs() {
	flag.Parse()
	if conf.Redis == "" && os.Getenv("REDIS") != "" {
		conf.Redis = os.Getenv("REDIS")
	} else if conf.Redis == "" && os.Getenv("REDIS") == "" {
		conf.Redis = defaultRedis
	}
	if conf.Password == "" && os.Getenv("PASSWORD") != "" {
		conf.Password = os.Getenv("PASSWORD")
	}
	if conf.Source == "" && os.Getenv("SOURCE") != "" {
		conf.Source = os.Getenv("SOURCE")
	} else if conf.Source == "" && os.Getenv("SOURCE") == "" {
		conf.Source = defaultSource
	}
	if conf.Db == 0 && os.Getenv("DB") != "" {
		var err error
		conf.Db, err = strconv.Atoi(os.Getenv("DB"))
		if err != nil {
			log.Fatalf("Error params %s", err)
		}
	}
	if conf.Port == "" && os.Getenv("PORT") != "" {
		conf.Port = os.Getenv("PORT")
	}
	if conf.Port == "" && os.Getenv("PORT") == "" {
		conf.Port = defaultPort
	}
}

// GetConf return configuration application
func GetConf() *Config {
	if *conf == (Config{}) {
		parseArgs()
	}
	return conf
}
