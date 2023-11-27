package filters

import (
	"github.com/shinhagunn/mpa/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func applyFilter(filters ...mongodb.Filter) []bson.M {
	value := []bson.M{}
	for _, filter := range filters {
		k, v := filter()
		value = append(value, bson.M{k: v})
	}

	return value
}

func WithAnd(filters ...mongodb.Filter) mongodb.Filter {
	return func() (k string, v interface{}) {
		return "$and", applyFilter(filters...)
	}
}

func WithOr(filters ...mongodb.Filter) mongodb.Filter {
	return func() (k string, v interface{}) {
		return "$or", applyFilter(filters...)
	}
}

func WithNot(filters ...mongodb.Filter) mongodb.Filter {
	return func() (k string, v interface{}) {
		return "$not", applyFilter(filters...)
	}
}

func WithNor(filters ...mongodb.Filter) mongodb.Filter {
	return func() (k string, v interface{}) {
		return "$nor", applyFilter(filters...)
	}
}
