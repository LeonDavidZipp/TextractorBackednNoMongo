db.createUser(
    {
        user: "imagedb_user",
        pwd: "70hHf01ghiuinaZgP",
        roles: [
        {
            role: "readWrite",
            db: "imagedb"
        }
        ]
    }
);