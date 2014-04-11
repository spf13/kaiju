package models

import (
	"time"

	"labix.org/v2/mgo/bson"
)

type CommentUser struct {
	Id       bson.ObjectId
	FullName string `bson:",omitempty"`
	Email    string
}

type Comment struct {
	Id        bson.ObjectId `bson:"_id"`
	User      CommentUser
	Forum     bson.ObjectId
	Timestamp time.Time
	Page      string
	Body      string
	Parent    *bson.ObjectId
	Ancestors []bson.ObjectId
}
