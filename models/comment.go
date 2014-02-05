package models

import (
    "labix.org/v2/mgo/bson"
)

type CommentUser struct {
    FullName string `bson:",omitempty"`
    UserId bson.ObjectId `bson:",omitempty"`
}

type Comment struct {
    Id bson.ObjectId `bson:"_id"`
    Page string `bson:",omitempty"`
    User CommentUser `bson:",omitempty"`
    Body string `bson:",omitempty"`
    Parent bson.ObjectId `bson:",omitempty"`
    Children []bson.ObjectId `bson:",omitempty"`
    Forum bson.ObjectId `bson:",omitempty"`
}
