package models

import (
    "labix.org/v2/mgo/bson"
)

type CommentUser struct {
    FullName string
    UserId bson.ObjectId
}

type Comment struct {
    Id bson.ObjectId "_id"
    Page string
    User CommentUser
    Body string
    Parent bson.ObjectId
    Children []bson.ObjectId
    Forum bson.ObjectId
}
