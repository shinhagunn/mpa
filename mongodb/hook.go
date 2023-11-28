package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type CallbackType string

var (
	BeforeCreate = CallbackType("before-create")
	AfterCreate  = CallbackType("after-create")
	BeforeUpdate = CallbackType("before-update")
	AfterUpdate  = CallbackType("after-update")
	BeforeDelete = CallbackType("before-delete")
	AfterDelete  = CallbackType("after-delete")
)

type action struct {
	before func()
	after  func()
}

type callBack struct {
	create action
	update action
	delete action
}

type collection struct {
	tabler   Tabler
	callback callBack
}

type Mongo struct {
	DB          *mongo.Database
	collections map[string]*collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		DB:          db,
		collections: make(map[string]*collection),
	}
}

func (d *Mongo) NewCollection(name string, tabler Tabler) {
	d.collections[name] = &collection{
		tabler: tabler,
		callback: callBack{
			create: action{
				before: func() {},
				after:  func() {},
			},
			update: action{
				before: func() {},
				after:  func() {},
			},
			delete: action{
				before: func() {},
				after:  func() {},
			},
		},
	}
}

func (d *Mongo) RegisterCallback(name string, callbackType CallbackType, fn func()) {
	switch callbackType {
	case BeforeCreate:
		d.collections[name].callback.create.before = fn
	case AfterCreate:
		d.collections[name].callback.create.after = fn
	case BeforeUpdate:
		d.collections[name].callback.update.before = fn
	case AfterUpdate:
		d.collections[name].callback.update.after = fn
	case BeforeDelete:
		d.collections[name].callback.delete.before = fn
	case AfterDelete:
		d.collections[name].callback.delete.after = fn
	}
}

func (d *Mongo) RunCallback(name string, callbackType CallbackType) {
	switch callbackType {
	case BeforeCreate:
		d.collections[name].callback.create.before()
	case AfterCreate:
		d.collections[name].callback.create.after()
	case BeforeUpdate:
		d.collections[name].callback.update.before()
	case AfterUpdate:
		d.collections[name].callback.update.after()
	case BeforeDelete:
		d.collections[name].callback.delete.before()
	case AfterDelete:
		d.collections[name].callback.delete.after()
	}
}
