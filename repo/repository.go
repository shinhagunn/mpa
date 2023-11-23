package repo

import (
	"context"

	filtersPkg "github.com/shinhagunn/mpa/filters"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r Repository) Count(ctx context.Context, filters ...filtersPkg.Filter) (int, error) {
	result, err := r.DB.Collection(r.Tabler.TableName()).CountDocuments(ctx, filtersPkg.ApplyFilters(filters...))
	return int(result), err
}

func (r Repository) Find(ctx context.Context, models interface{}, filters []filtersPkg.Filter, opts ...*options.FindOptions) error {
	cursor, err := r.DB.Collection(r.Tabler.TableName()).Find(ctx, filtersPkg.ApplyFilters(filters...), opts...)
	if err != nil {
		return err
	}

	if err := cursor.All(ctx, models); err != nil {
		return err
	}

	return nil
}

func (r Repository) First(ctx context.Context, model interface{}, filters ...filtersPkg.Filter) error {
	if err := r.DB.Collection(r.Tabler.TableName()).FindOne(ctx, filtersPkg.ApplyFilters(filters...)).Decode(model); err != nil {
		return err
	}

	return nil
}

func (r Repository) Last(ctx context.Context, model interface{}, filters ...filtersPkg.Filter) error {
	if err := r.DB.Collection(r.Tabler.TableName()).FindOne(
		ctx,
		filtersPkg.ApplyFilters(filters...),
		options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}),
	).Decode(model); err != nil {
		return err
	}

	return nil
}

func (r Repository) FirstOrCreate(ctx context.Context, model interface{}, filters ...filtersPkg.Filter) error {
	// _, err := r.DB.Collection(r.Tabler.TableName()).F(ctx, model)
	return nil
}

func (r Repository) Create(ctx context.Context, model interface{}, filters ...filtersPkg.Filter) error {
	_, err := r.DB.Collection(r.Tabler.TableName()).InsertOne(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) UpdateByID(ctx context.Context, id primitive.ObjectID, value interface{}, filters ...filtersPkg.Filter) error {
	updateValue := make(map[string]interface{})
	updateValue["$set"] = value

	if _, err := r.DB.Collection(r.Tabler.TableName()).UpdateByID(ctx, id, updateValue); err != nil {
		return err
	}

	return nil
}

func (r Repository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.DB.Collection(r.Tabler.TableName()).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
