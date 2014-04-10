db.getSiblingDB("kaiju").dropDatabase()

db.getSiblingDB("kaiju").users.insert({
    _id: new ObjectId("5346e476331583002c7de60d"),
    fullname: "John Doe",
    email: "fake_email@gmail.com",
    website: "http://example.com",
    location: "NYC",
    bio: "...",
    services: [
        { type: "facebook", user_id: '1234...' },
        { type: "twitter", user_id: '1234...' },
        { type: "google", user_id: "example@gmail.com" }
    ]
})

db.getSiblingDB("kaiju").users.insert({
    _id: new ObjectId("5346e476331583002c7de60e"),
    fullname: "Jane Doe",
    email: "different_fake_email@gmail.com",
    website: "http://example.com",
    location: "NYC",
    bio: "...",
    services: [
    ]
})

var forum = {
    _id: new ObjectId("5346e494331583002c7de60e"),
    name: "Steve's blog on how to improve soccer",
    adminusers: [new ObjectId(), new ObjectId()]
}
db.getSiblingDB("kaiju").forums.insert(forum)

db.getSiblingDB("kaiju").comments.insert([{ "_id" : ObjectId("5346f749080fcc8830000004"), "user" : { "fullname" : "John Doe", "userid" : ObjectId("5346e476331583002c7de60d"), "email" : "fake_email@gmail.com" }, "forum" : ObjectId("5346e494331583002c7de60e"), "page" : "dans sc2 blog post", "body" : "sc2 is awesome. Toss OP", "parent" : null, "ancestors" : [ ] },
{ "_id" : ObjectId("5346f74a080fcc8830000005"), "user" : { "fullname" : "John Doe", "userid" : ObjectId("5346e476331583002c7de60d"), "email" : "fake_email@gmail.com" }, "forum" : ObjectId("5346e494331583002c7de60e"), "page" : "dans sc2 blog post", "body" : "sc2 is awesome. Toss OP", "parent" : null, "ancestors" : [ ] },
{ "_id" : ObjectId("5346f7a8080fcc885c000001"), "user" : { "id" : ObjectId("5346e476331583002c7de60d"), "fullname" : "John Doe", "email" : "fake_email@gmail.com" }, "forum" : ObjectId("5346e494331583002c7de60e"), "page" : "dans sc2 blog post", "body" : "sc2 is awesome. Toss OP", "parent" : ObjectId("5346f749080fcc8830000004"), "ancestors" : [ ObjectId("5346f749080fcc8830000004") ] },
{ "_id" : ObjectId("5346f7d1080fcc886c000001"), "user" : { "id" : ObjectId("5346e476331583002c7de60d"), "fullname" : "John Doe", "email" : "fake_email@gmail.com" }, "forum" : ObjectId("5346e494331583002c7de60e"), "page" : "dans sc2 blog post", "body" : "sc2 is awesome. Toss OP", "parent" : ObjectId("5346f7a8080fcc885c000001"), "ancestors" : [ ObjectId("5346f749080fcc8830000004"), ObjectId("5346f7a8080fcc885c000001") ] }])
