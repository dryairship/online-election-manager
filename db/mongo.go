package db

import (
    "github.com/dryairship/online-election-manager/config"
    "gopkg.in/mgo.v2"
)

// Struct that is used to store the session of the database connection.
type ElectionDatabase struct {
    Session *mgo.Session
}

// Function to establish database connection.
func ConnectToDatabase() (ElectionDatabase, error) {
    sess, err := mgo.Dial(config.MongoDialURL)
    if err != nil {
        return ElectionDatabase{nil}, err
    }
    err = sess.DB(config.MongoDbName).Login(config.MongoUsername, config.MongoPassword)
    return ElectionDatabase{sess}, err
}

// Function to delete all entries from the database.
func (db ElectionDatabase) ResetDatabase() error {
    _, err := db.Session.DB(config.MongoDbName).C("candidates").RemoveAll(nil)
    if err != nil {
        return err
    }
    
    _, err = db.Session.DB(config.MongoDbName).C("voters").RemoveAll(nil)
    if err != nil {
        return err
    }
    
    _, err = db.Session.DB(config.MongoDbName).C("votes").RemoveAll(nil)
    if err != nil {
        return err
    }
    
    _, err = db.Session.DB(config.MongoDbName).C("posts").RemoveAll(nil)
    if err != nil {
        return err
    }
    
    _, err = db.Session.DB(config.MongoDbName).C("ceo").RemoveAll(nil)
    return err
}

