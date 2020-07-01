package mongo

import (
	"context"
	"time"

	"github.com/divyanshgaba/bingo/config"
	kitlog "github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	timeout  = 60 * time.Second
	database string
)

func init() {
	database = config.Mongo().Database
}

// NewClient creates a new mongo client.
func NewClient(logger kitlog.Logger) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Mongo().URI()).SetReadPreference(readpref.Primary()))
	if err != nil {
		logger.Log("err", err, "msg", "mongo: error while creating mongo client")
		return client, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Log("err", err, "msg", "mongo: error while connecting with mongo")
		return client, err
	}
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Log("err", err, "msg", "mongo: ping failed to mongo")
		return client, err
	}
	logger.Log("msg", "mongo: connection successful")
	return client, err
}

// // NewClient creates a new mongo client.
// func NewClient() (*mongo.Client, error) {
// 	client, err := mongo.NewClient(options.Client().ApplyURI(config.Mongo().URI()).SetReadPreference(readpref.Primary()))
// 	if err != nil {
// 		log.WithField("err", err).Error("mongoclient: error while creating mongo client")
// 		return client, err
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.WithField("err", err).Error("mongoclient: error while connecting with mongo")
// 		return client, err
// 	}
// 	ctx, cancel = context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	err = client.Ping(ctx, readpref.Primary())
// 	if err != nil {
// 		log.WithField("err", err).Error("mongoclient: ping failed to mongo")
// 		return client, err
// 	}
// 	return client, err
// }
