package db

import(
    "github.com/dryairship/online-election-manager/models"
    "github.com/dryairship/online-election-manager/config"
    "gopkg.in/mgo.v2/bson"
)

// Function to find the candidate with that username in the database.
func (db ElectionDatabase) GetCandidate(username string) (models.Candidate, error) {
    candidatesCollection := db.Session.DB(config.MongoDbName).C("candidates")
    candidate := models.Candidate{}
    err := candidatesCollection.Find(bson.M{"username":username}).One(&candidate)
    return candidate, err
}

// Function to find all candidates in the database.
func (db ElectionDatabase) GetAllCandidates() ([]models.Candidate, error) {
    candidatesCollection := db.Session.DB(config.MongoDbName).C("candidates")
    var candidates []models.Candidate
    err := candidatesCollection.Find(nil).All(&candidates)
    return candidates, err
}

// Function to update the details of a candidate.
func (db ElectionDatabase) UpdateCandidate(username string, newCandidate *models.Candidate) error {
    candidatesCollection := db.Session.DB(config.MongoDbName).C("candidates")
    err := candidatesCollection.Update(bson.M{"username":username}, &newCandidate)
    return err
}

// Function to insert a new candidate in the database.
func (db ElectionDatabase) AddNewCandidate(candidate *models.Candidate) error {
    candidatesCollection := db.Session.DB(config.MongoDbName).C("candidates")
    err := candidatesCollection.Insert(&candidate)
    return err
}
