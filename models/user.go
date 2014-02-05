package models

import {
    "labix.org/v2/mgo/bson"
}

type UserService struct {
    Provider string
    UserId string
}

type User struct {
    Id bson.ObjectId "_id"
    FullName string
    Website string
    Location string
    Bio string
    Services []UserService
    Forum bson.ObjectId
}
