db.createUser(
    {
        user: "appUser",
        pwd: "appUser123",
        roles:[
            {
                role : "readWrite",
                db : "bingo"
            }
        ]
    }
)