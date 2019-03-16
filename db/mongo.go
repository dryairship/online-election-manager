package db

import (
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/models"
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

func (db ElectionDatabase) FindStudentSkeleton(roll string) (models.StudentSkeleton, error) {
    studentsCollection := db.Session.DB(config.MongoDbName).C("students")
    skeleton := models.StudentSkeleton{}
    err := studentsCollection.Find(bson.M{"roll":roll}).One(&skeleton)
    return skeleton, err
}

func (db ElectionDatabase) FindVoter(roll string) (models.Voter, error) {
    votersCollection := db.Session.DB(config.MongoDbName).C("voters")
    voter := models.Voter{}
    err := votersCollection.Find(bson.M{"roll":roll}).One(&voter)
    return voter, err
}

func (db ElectionDatabase) AddNewVoter(voter *models.Voter) error {
    votersCollection := db.Session.DB(config.MongoDbName).C("voters")
    err := votersCollection.Insert(&voter)
    return err
}
