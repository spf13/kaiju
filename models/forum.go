package models

import (
    "labix.org/v2/mgo/bson"
)

type Forum struct {
    Id bson.ObjectId "_id"
    Name string
    AdminUsers []bson.ObjectId
}
