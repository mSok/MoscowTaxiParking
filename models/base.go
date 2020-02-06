package models

import (
	"app/utils"
	"log"

	"github.com/go-redis/redis"
)

var db *DBClient

// DBClient database client
type DBClient struct {
	db *redis.Client
}

// Create DB connection
func configure() {

	var conf = utils.GetConf()

	dbredis := redis.NewClient(&redis.Options{
		Addr:     conf.Redis,
		Password: conf.Password,
		DB:       conf.Db,
	})
	log.Printf("Connected to DB Address %s db: %d", conf.Redis, conf.Db)
	db = &DBClient{db: dbredis}
}

// GetDB return db singleton
func GetDB() *DBClient {
	if db == nil {
		configure()
	}
	return db
}
