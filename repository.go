package main

import (
	"context"

	"github.com/shinhagunn/mpa/filters"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tabler interface {
	TableName() string
}

type Repository struct {
	DB     *mongo.Database
	Tabler Tabler
}

func New(db *mongo.Database, entity Tabler) Repository {
	return Repository{db, entity}
}

func (r Repository) Count(ctx context.Context, filters ...filters.Filter) (int, error) {
	result, err := r.DB.Collection(r.Tabler.TableName()).CountDocuments(ctx, ApplyFilters(filters...))
	return int(result), err
}

func (r Repository) Find(ctx context.Context, models interface{}, filters []filters.Filter, opts *options.FindOptions) error {
	cursor, err := r.DB.Collection(r.Tabler.TableName()).Find(ctx, ApplyFilters(filters...), opts)
	if err != nil {
		return err
	}

	if err := cursor.All(ctx, &models); err != nil {
		return err
	}

	return nil
}

func (r Repository) First(ctx context.Context, model interface{}, filters []filters.Filter) error {
	result := r.DB.Collection(r.Tabler.TableName()).FindOne(ctx, ApplyFilters(filters...))

	if err := result.Decode(&model); err != nil {
		return err
	}

	return nil
}

func (r Repository) Last(ctx context.Context, model interface{}, filters []filters.Filter) error {
	result := r.DB.Collection(r.Tabler.TableName()).FindOne(
		ctx,
		ApplyFilters(filters...),
		options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}),
	)

	if err := result.Decode(&model); err != nil {
		return err
	}

	return nil
}

func (r Repository) FirstOrCreate(ctx context.Context, model interface{}, filters []filters.Filter) error {
	return nil
}

func (r Repository) Create(ctx context.Context, model interface{}, filters []filters.Filter) error {
	_, err := r.DB.Collection(r.Tabler.TableName()).InsertOne(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) Updates(ctx context.Context, model interface{}, filters []filters.Filter) error {
	_, err := r.DB.Collection(r.Tabler.TableName()).UpdateOne(ctx, ApplyFilters(filters...), model)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) Delete(ctx context.Context, model interface{}, filters []filters.Filter) error {
	_, err := r.DB.Collection(r.Tabler.TableName()).DeleteOne(ctx, ApplyFilters(filters...))
	if err != nil {
		return err
	}

	return nil
}
