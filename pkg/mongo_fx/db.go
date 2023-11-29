package mongo_fx

import (
	"context"
	"fmt"
	"os"

	"github.com/shinhagunn/mpa/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(cfg *config.Config) (*Mongo, error) {
	myLogger := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			PrettyPrint:     true,
		},
	}

	sink := &DBlogger{log: myLogger}

	loggerOptions := options.
		Logger().
		SetSink(sink).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	options := options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d",
			cfg.Database.User,
			cfg.Database.Pass,
			cfg.Database.Host,
			cfg.Database.Port,
		)).
		SetLoggerOptions(loggerOptions)

	client, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		return nil, err
	}

	return NewMongo(client.Database(cfg.Database.Name)), nil
}
