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

func (db ElectionDatabase) UpdateVoter(roll string, newVoter *models.Voter) error {
    votersCollection := db.Session.DB(config.MongoDbName).C("voters")
    err := votersCollection.Update(bson.M{"roll":roll},&newVoter)
    return err
}

func (db ElectionDatabase) GetPosts() ([]models.Post, error){
    postsCollection := db.Session.DB(config.MongoDbName).C("posts")
    posts := []models.Post{}
    err := postsCollection.Find(nil).All(&posts)
    return posts, err
}

func (db ElectionDatabase) GetCandidate(username string) (models.Candidate, error) {
    candidatesCollection := db.Session.DB(config.MongoDbName).C("candidates")
    candidate := models.Candidate{}
    err := candidatesCollection.Find(bson.M{"username":username}).One(&candidate)
    return candidate, err
}

func (db ElectionDatabase) InsertVote(vote *models.Vote) error {
    votesCollection := db.Session.DB(config.MongoDbName).C("votes")
    err := votesCollection.Insert(&vote)
    return err
}

func (db ElectionDatabase) GetCEO() (models.CEO, error) {
    ceoCollection := db.Session.DB(config.MongoDbName).C("ceo")
    ceo := models.CEO{}
    err := ceoCollection.Find(nil).One(&ceo)
    return ceo, err
}

func (db ElectionDatabase) InsertCEO(ceo *models.CEO) error {
    ceoCollection := db.Session.DB(config.MongoDbName).C("ceo")
    err := ceoCollection.Insert(&ceo)
    return err
}

func (db ElectionDatabase) UpdateCEO(newCEO *models.CEO) error {
    ceoCollection := db.Session.DB(config.MongoDbName).C("ceo")
    err := ceoCollection.Update(bson.M{"username":"CEO"}, &newCEO)
    return err
}

func (db ElectionDatabase) GetVotes() ([]models.Vote, error) {
    votesCollection := db.Session.DB(config.MongoDbName).C("votes")
    var votes []models.Vote
    err := votesCollection.Find(nil).All(&votes)
    return votes, err
}

