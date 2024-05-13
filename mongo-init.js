db.createUser(
    {
        user: "imagedb_user",
        pwd: "test1234",
        roles: [
        {
            role: "readWrite",
            db: "imagedb"
        }
        ]
    }
);