package mpa_fx

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

type Repository[T Tabler] interface {
	DB() *mongo.Database
	TableName() string
	Transaction(ctx context.Context, handler func(sessionContext mongo.SessionContext) error) error
	Count(ctx context.Context, filters ...filtersPkg.Filter) (int, error)
	Find(ctx context.Context, filters []filtersPkg.Filter, opts ...*options.FindOptions) (models []*T, err error)
	First(ctx context.Context, filters ...filtersPkg.Filter) (model *T, err error)
	Last(ctx context.Context, filters ...filtersPkg.Filter) (model *T, err error)
	FirstOrCreate(ctx context.Context, create *T, filters ...filtersPkg.Filter) (*T, error)
	Create(ctx context.Context, model *T, filters ...filtersPkg.Filter) error
	Updates(ctx context.Context, model *T, value map[string]interface{}, filters ...filtersPkg.Filter) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type repository[T Tabler] struct {
	db     *mongo.Database
	tabler Tabler
}

func NewRepository[T Tabler](db *mongo.Database, entity T) Repository[T] {
	return repository[T]{db, entity}
}

func (r repository[T]) DB() *mongo.Database {
	return r.db
}

func (r repository[T]) TableName() string {
	return r.tabler.TableName()
}

func (r repository[T]) Transaction(ctx context.Context, handler func(sessionContext mongo.SessionContext) error) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	if err := mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if err := handler(sessionContext); err != nil {
			return err
		}

		if err := session.CommitTransaction(sessionContext); err != nil {
			return err
		}

		return nil
	}); err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			return abortErr
		}

		return err
	}

	return nil
}

func (r repository[T]) Count(ctx context.Context, filters ...filtersPkg.Filter) (int, error) {
	result, err := r.db.Collection(r.tabler.TableName()).CountDocuments(ctx, filtersPkg.ApplyFilters(filters...))
	return int(result), err
}

func (r repository[T]) Find(ctx context.Context, filters []filtersPkg.Filter, opts ...*options.FindOptions) (models []*T, err error) {
	cursor, err := r.db.Collection(r.tabler.TableName()).Find(ctx, filtersPkg.ApplyFilters(filters...), opts...)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	return models, nil
}

func (r repository[T]) First(ctx context.Context, filters ...filtersPkg.Filter) (model *T, err error) {
	if err := r.db.Collection(r.tabler.TableName()).FindOne(ctx, filtersPkg.ApplyFilters(filters...)).Decode(&model); err != nil {
		return nil, err
	}

	return model, nil
}

func (r repository[T]) Last(ctx context.Context, filters ...filtersPkg.Filter) (model *T, err error) {
	if err := r.db.Collection(r.tabler.TableName()).FindOne(
		ctx,
		filtersPkg.ApplyFilters(filters...),
		options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}}),
	).Decode(&model); err != nil {
		return nil, err
	}

	return model, nil
}

func (r repository[T]) FirstOrCreate(ctx context.Context, create *T, filters ...filtersPkg.Filter) (*T, error) {
	var model *T
	err := r.db.Collection(r.tabler.TableName()).FindOne(ctx, filtersPkg.ApplyFilters(filters...)).Decode(&model)
	if err == nil {
		return model, nil
	}

	if _, err := r.db.Collection(r.tabler.TableName()).InsertOne(ctx, create); err != nil {
		return nil, err
	}

	return create, nil
}

func (r repository[T]) Create(ctx context.Context, model *T, filters ...filtersPkg.Filter) error {
	if _, err := r.db.Collection(r.tabler.TableName()).InsertOne(ctx, model); err != nil {
		return err
	}

	return nil
}

// Value must be map[string]interface{} because if is *T, if struct have a field not fill value, mongoDB will update that field zero value
func (r repository[T]) Updates(ctx context.Context, model *T, value map[string]interface{}, filters ...filtersPkg.Filter) error {
	result, err := r.db.Collection(r.tabler.TableName()).UpdateOne(ctx, &model, map[string]interface{}{"$set": value})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r repository[T]) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.Collection(r.tabler.TableName()).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
