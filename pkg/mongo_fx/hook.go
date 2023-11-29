package mongo_fx

import "go.mongodb.org/mongo-driver/mongo"

type CallbackType string

var (
	CallbackTypeBeforeCreate = CallbackType("before-create")
	CallbackTypeAfterCreate  = CallbackType("after-create")
	CallbackTypeBeforeUpdate = CallbackType("before-update")
	CallbackTypeAfterUpdate  = CallbackType("after-update")
	CallbackTypeBeforeDelete = CallbackType("before-delete")
	CallbackTypeAfterDelete  = CallbackType("after-delete")
)

type CallbackFunc func(db *mongo.Database, value interface{}) error

type collection struct {
	tabler   Tabler
	callback map[CallbackType][]CallbackFunc
}

type Mongo struct {
	DB          *mongo.Database
	Collections map[string]*collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		DB:          db,
		Collections: make(map[string]*collection),
	}
}
