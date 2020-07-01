db.createUser(
    {
        user: "testUser",
        pwd: "testUser123",
        roles:[
            {
                role : "readWrite",
                db : "bingo-test"
            }
        ]
    }
)