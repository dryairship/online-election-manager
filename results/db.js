const mongoose = require('mongoose');

class DBConnection {
    static dial() {
        if (this.db)
            return Promise.resolve(this.db);

        mongoose.connect(process.env.OEMMongoDialURL, this.options);
        this.db = mongoose.connection;

        this.db.on('error', console.error.bind(console, '[ERROR] [MongoDB]'));
    }
}
DBConnection.db = null
DBConnection.options = {
    auth: {
        authSource: process.env.OEMMongoDbName
    },
    dbName: process.env.OEMMongoDbName,
    user: process.env.OEMMongoUsername,
    pass: process.env.OEMMongoPassword,
    useNewUrlParser: true
}

module.exports = DBConnection;