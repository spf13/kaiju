db.getSiblingDB("kaiju").dropDatabase()

var user = {
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
}
db.getSiblingDB("kaiju").users.insert(user)

var forum = {
    _id: new ObjectId("5346e494331583002c7de60e"),
    name: "Steve's blog on how to improve soccer",
    adminusers: [new ObjectId(), new ObjectId()]
}
db.getSiblingDB("kaiju").forums.insert(forum)
