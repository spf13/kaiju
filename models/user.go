package models

import (
    "labix.org/v2/mgo/bson"
)

type UserService struct {
    Provider string
    UserId string
}

type User struct {
    Id bson.ObjectId `bson:"_id"`
    FullName string `bson:",omitempty"`
	Email string `bson:",omitempty"`
    Website string `bson:",omitempty"`
    Location string `bson:",omitempty"`
    Bio string `bson:",omitempty"`
    Services []UserService `bson:",omitempty"`
    Forum bson.ObjectId `bson:",omitempty"`
}
