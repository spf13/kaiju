db.getSiblingDB("kaiju").dropDatabase()

var user = {
    _id: new ObjectId("5346e476331583002c7de60d"),
    FullName: "John Doe",
    Email: "fake_email@gmail.com",
    Website: "http://example.com",
    Location: "NYC",
    Bio: "...",
    Services: [
        { type: "facebook", user_id: '1234...' },
        { type: "twitter", user_id: '1234...' },
        { type: "google", user_id: "example@gmail.com" }
    ]
}
db.getSiblingDB("kaiju").users.insert(user)

var commentA = {
    _id: new ObjectId(),
    Page: 'http://example.com/article/1',
    User: {
        FullName: "John Doe",
        Email: "jdoe@example.com",
    },
    Body: '<p>foo</p>',
    Parent: null,
    metadata: {
        ip: '127.0.0.1'
    }
};

var commentAA = {
    _id: new ObjectId(),
    Page: 'http://example.com/article/1',
    User: {
        FullName: "Mary Doe",
        Email: "mdoe@example.com",
    },
    Body: '<p>bar</p>',
    Parent: commentA._id,
    metadata: {
        ip: '127.0.0.1'
    }
};

commentA.children.push(commentAA._id);
