package mongo_fx

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/shinhagunn/mpa/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tabler interface {
	TableName() string
}

type callback[T Tabler] map[CallbackType][]func(db *mongo.Database, value *T) error

type Repository[T Tabler] interface {
	DB() *mongo.Database
	TableName() string
	AddCallback(kind CallbackType, callback func(db *mongo.Database, value *T) error)
	Transaction(ctx context.Context, handler func(sessionContext mongo.SessionContext) error) error
	Count(ctx context.Context, filters ...mongodb.Filter) (int, error)
	Find(ctx context.Context, opts *options.FindOptions, filters ...mongodb.Filter) (models []*T, err error)
	First(ctx context.Context, filters ...mongodb.Filter) (model *T, err error)
	Last(ctx context.Context, filters ...mongodb.Filter) (model *T, err error)
	FirstOrCreate(ctx context.Context, model *T, create *T, filters ...mongodb.Filter) error
	Create(ctx context.Context, model *T, filters ...mongodb.Filter) error
	Updates(ctx context.Context, model *T, value map[string]interface{}, filters ...mongodb.Filter) error
	Delete(ctx context.Context, filters ...mongodb.Filter) error
}

type repository[T Tabler] struct {
	callback   callback[T]
	repository mongodb.Repository
	entity     T
}

func NewRepository[T Tabler](m *Mongo, entity T) Repository[T] {
	return repository[T]{
		callback:   make(callback[T]),
		repository: mongodb.New(m.DB, entity),
	}
}

func (r repository[T]) DB() *mongo.Database {
	return r.repository.DB
}

func (r repository[T]) TableName() string {
	return r.entity.TableName()
}

func (r repository[T]) AddCallback(kind CallbackType, callback func(db *mongo.Database, value *T) error) {
	if len(r.callback[kind]) == 0 {
		r.callback[kind] = []func(db *mongo.Database, value *T) error{callback}
	} else {
		r.callback[kind] = append(r.callback[kind], callback)
	}
}

func (r repository[T]) Transaction(ctx context.Context, handler func(sessionContext mongo.SessionContext) error) error {
	session, err := r.repository.DB.Client().StartSession()
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

func (r repository[T]) Count(ctx context.Context, filters ...mongodb.Filter) (int, error) {
	count, err := r.repository.Count(ctx, filters...)
	if err != nil {
		return 0, errors.Wrap(err, "repository count")
	}

	return count, err
}

func (r repository[T]) Find(ctx context.Context, opts *options.FindOptions, filters ...mongodb.Filter) (models []*T, err error) {
	err = r.repository.Find(ctx, &models, opts, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository find")
	}

	return models, nil
}

func (r repository[T]) First(ctx context.Context, filters ...mongodb.Filter) (model *T, err error) {
	err = r.repository.First(ctx, &model, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository first")
	}

	return model, nil
}

func (r repository[T]) Last(ctx context.Context, filters ...mongodb.Filter) (model *T, err error) {
	err = r.repository.Last(ctx, &model, filters...)
	if err != nil {
		return nil, errors.Wrap(err, "repository last")
	}

	return model, nil
}

func (r repository[T]) FirstOrCreate(ctx context.Context, model *T, create *T, filters ...mongodb.Filter) error {
	err := r.repository.FirstOrCreate(ctx, model, &create, filters...)
	if err != nil {
		return errors.Wrap(err, "repository first or create")
	}

	return nil
}

func (r repository[T]) Create(ctx context.Context, model *T, filters ...mongodb.Filter) error {
	if err := r.runCallback(CallbackTypeBeforeCreate, model); err != nil {
		return err
	}

	err := r.repository.Create(ctx, &model, filters...)
	if err != nil {
		return errors.Wrap(err, "repository create")
	}

	if err := r.runCallback(CallbackTypeAfterCreate, model); err != nil {
		return err
	}

	return nil
}

func (r repository[T]) Updates(ctx context.Context, model *T, value map[string]interface{}, filters ...mongodb.Filter) error {
	if err := r.runCallback(CallbackTypeBeforeUpdate, model); err != nil {
		return err
	}

	if err := r.repository.Updates(ctx, &model, value, filters...); err != nil {
		return errors.Wrap(err, "repository update")
	}

	if err := r.updateModel(model, value); err != nil {
		return err
	}

	if err := r.runCallback(CallbackTypeAfterUpdate, model); err != nil {
		return err
	}

	return nil
}

func (r repository[T]) Delete(ctx context.Context, filters ...mongodb.Filter) error {
	err := r.repository.Delete(ctx, filters...)
	if err != nil {
		return errors.Wrap(err, "repository delete")
	}

	return nil
}

func (r repository[T]) updateModel(model *T, update map[string]interface{}) error {
	objValue := reflect.ValueOf(model)

	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	for fieldName, newValue := range update {
		arr := strings.Split(fieldName, "_")

		s := ""
		for _, item := range arr {
			s += strings.ToUpper(string(item[0])) + item[1:]
		}

		field := objValue.FieldByName(s)

		if field.IsValid() {
			newFieldValue := reflect.ValueOf(newValue)
			if newFieldValue.Type().ConvertibleTo(field.Type()) {
				field.Set(newFieldValue.Convert(field.Type()))
			} else {
				return fmt.Errorf("cannot convert %v to %v", newFieldValue.Type(), field.Type())
			}
		}
	}

	return nil
}

func (r repository[T]) runCallback(callbackType CallbackType, model *T) error {
	for _, callback := range r.callback[callbackType] {
		if err := callback(r.repository.DB, model); err != nil {
			return err
		}
	}

	return nil
}
