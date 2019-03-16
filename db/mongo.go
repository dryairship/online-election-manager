package db

import (
    "github.com/dryairship/online-election-manager/config"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type ElectionDatabase struct {
    Session *mgo.Session
}


func ConnectToDatabase() (ElectionDatabase, error) {
    sess, err := mgo.Dial(config.MongoDialURL)
    if err != nil {
        return ElectionDatabase{nil}, err
    }
    err = sess.DB(config.MongoDbName).Login(config.MongoUsername, config.MongoPassword)
    return ElectionDatabase{sess}, err
}

