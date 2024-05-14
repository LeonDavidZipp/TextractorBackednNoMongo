db.createUser(
    {
        user: "imagedb_user",
        pwd: "test123",
        roles: [
        {
            role: "readWrite",
            db: "imagedb"
        }
        ]
    }
);