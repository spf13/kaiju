package models

import (
    "labix.org/v2/mgo/bson"
)

type Forum struct {
    Id bson.ObjectId `bson:"_id"`
    Name string `bson:",omitempty"`
    AdminUsers []bson.ObjectId `bson:",omitempty"`
}
