var user = {
    _id: new ObjectId(),
    fullName: "John Doe",
    website: "http://example.com",
    location: "NYC",
    bio: "...",
    services: [
        { type: "facebook", user_id: '1234...' },
        { type: "twitter", user_id: '1234...' },
        { type: "google", user_id: "example@gmail.com" }
    ]
}

var commentA = {
    _id: new ObjectId(),
    page: 'http://example.com/article/1',
    user: {
        name: "John Doe",
        email: "jdoe@example.com",
    },
    body: '<p>foo</p>',
    parent: null,
    children: [],
    metadata: {
        ip: '127.0.0.1'
    }
};

var commentAA = {
    _id: new ObjectId(),
    page: 'http://example.com/article/1',
    user: {
        name: "Mary Doe",
        email: "mdoe@example.com",
    },
    body: '<p>bar</p>',
    parent: commentA._id,
    children: [],
    metadata: {
        ip: '127.0.0.1'
    }
};

commentA.children.push(commentAA._id);

