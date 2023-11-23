package filters

import "go.mongodb.org/mongo-driver/bson"

type Filter func() bson.E
