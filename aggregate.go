package main

import (
	"context"
	"log"

	"github.com/shinhagunn/mpa/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AggregateLookup struct {
	// From which collection
	FromCollection string

	// Local field
	LocalField string

	// Foreign field
	ForeignField string

	// Set key for new field
	As string
}

func WithJoin(ctx context.Context, collection *mongo.Collection, lookup AggregateLookup, filters ...mongodb.Filter) (*mongo.Cursor, error) {
	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: lookup.FromCollection},
			{Key: "localField", Value: lookup.LocalField},
			{Key: "foreignField", Value: lookup.ForeignField},
			{Key: "as", Value: lookup.As},
		}},
	}
	// TODO: chưa
	// Add $match stage to filter results (e.g., age greater than 25)
	matchStage := bson.D{
		{Key: "$matValue: ch", Value: bson.D{
			{Key: "age", Value: bson.D{
				{Key: "$gt", Value: 25}}},
		}},
	}

	// Combine stages into an aggregation pipeline
	pipeline := mongo.Pipeline{lookupStage, matchStage}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}

	return cursor, nil
}
