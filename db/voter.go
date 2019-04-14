package db

import(
    "github.com/dryairship/online-election-manager/models"
    "github.com/dryairship/online-election-manager/config"
    "gopkg.in/mgo.v2/bson"
)

// Function to find the basic details of a student from its roll number.
func (db ElectionDatabase) FindStudentSkeleton(roll string) (models.StudentSkeleton, error) {
    studentsCollection := db.Session.DB(config.MongoDbName).C("students")
    skeleton := models.StudentSkeleton{}
    err := studentsCollection.Find(bson.M{"roll":roll}).One(&skeleton)
    return skeleton, err
}

// Function to find a voter in the database.
func (db ElectionDatabase) FindVoter(roll string) (models.Voter, error) {
    votersCollection := db.Session.DB(config.MongoDbName).C("voters")
    voter := models.Voter{}
    err := votersCollection.Find(bson.M{"roll":roll}).One(&voter)
    return voter, err
}

// Function to insert a new voter into the database.
func (db ElectionDatabase) AddNewVoter(voter *models.Voter) error {
    votersCollection := db.Session.DB(config.MongoDbName).C("voters")
    err := votersCollection.Insert(&voter)
    return err
}

// Function to update the details of a voter.
func (db ElectionDatabase) UpdateVoter(roll string, newVoter *models.Voter) error {
    votersCollection := db.Session.DB(config.MongoDbName).C("voters")
    err := votersCollection.Update(bson.M{"roll":roll},&newVoter)
    return err
}
