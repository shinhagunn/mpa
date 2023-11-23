package main

import (
	"github.com/shinhagunn/mpa/filters"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ordering string

var (
	OrderingAsc  Ordering = "asc"
	OrderingDesc Ordering = "desc"
)

type OptionsFind struct {
	Page     int
	Limit    int
	OrderBy  string
	Ordering Ordering
}

func NewOptionsFind(optionsFind OptionsFind) *options.FindOptions {
	opts := options.Find()

	if optionsFind.Page > 0 {
		opts.SetLimit(int64(optionsFind.Page))
	}

	if optionsFind.Limit > 0 {
		opts.SetLimit(int64(optionsFind.Limit))
	}

	if len(optionsFind.OrderBy) > 0 && len(optionsFind.Ordering) > 0 {
		order := 1
		if optionsFind.Ordering == OrderingDesc {
			order = -1
		}

		opts.SetSort(bson.M{optionsFind.OrderBy: order})
	}

	return opts
}

func ApplyFilters(fs ...filters.Filter) bson.D {
	var result []bson.E

	for _, f := range fs {
		result = append(result, f())
	}

	return result
}
