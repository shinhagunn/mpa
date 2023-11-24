package main

import (
	"context"

	"github.com/shinhagunn/mpa/config"
	"github.com/shinhagunn/mpa/pkg/mpa_fx"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := mpa_fx.New(cfg)
	if err != nil {
		panic(err)
	}

	session, err := db.Client().StartSession()
	if err != nil {
		panic(err)
	}

	defer session.EndSession(context.TODO())

	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if err := session.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			panic(abortErr)
		}
		panic(err)
	}
}
