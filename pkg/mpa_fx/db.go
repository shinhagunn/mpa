package mpa_fx

import (
	"context"
	"fmt"

	"github.com/shinhagunn/mpa/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type mongoParams struct {
// 	fx.In

// 	Config config.DB
// }

func New(cfg *config.Config) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Host,
		cfg.Database.Port,
	))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	client.StartSession()

	return client.Database(cfg.Database.Name), nil
}
