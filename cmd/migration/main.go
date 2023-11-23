package main

import (
	"context"
	"fmt"

	"github.com/shinhagunn/mpa/config"
	"github.com/shinhagunn/mpa/models"
	"github.com/shinhagunn/mpa/pkg/mpa_fx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

type Indexer interface {
	TableName() string
	SetIndex() map[string]mpa_fx.IndexType
}

func GetOptionIndex(name string, typeIndex mpa_fx.IndexType) *options.IndexOptions {
	opts := options.Index().SetName(name)

	switch typeIndex {
	case mpa_fx.IndexTypeUnique:
		opts.SetUnique(true)
	default:
		opts.SetUnique(true)
	}

	return opts
}

func main() {
	app := fx.New(
		config.Module,
		mpa_fx.Module,
		fx.Invoke(func(db *mongo.Database) error {
			indexers := []Indexer{
				&models.User{},
			}

			for _, indexer := range indexers {
				for field, typeIndex := range indexer.SetIndex() {
					if _, err := db.Collection(indexer.TableName()).Indexes().CreateOne(
						context.TODO(),
						mongo.IndexModel{
							Keys:    bson.D{{Key: field, Value: 1}},
							Options: GetOptionIndex(fmt.Sprintf("index_%s_on_%s", indexer.TableName(), field), typeIndex),
						},
					); err != nil {
						return err
					}
				}
			}

			return nil
		}),
	)

	defer app.Stop(context.Background())

	app.Start(context.Background())
}
