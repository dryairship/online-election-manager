package db

import(
    "github.com/dryairship/online-election-manager/models"
    "github.com/dryairship/online-election-manager/config"
    "gopkg.in/mgo.v2/bson"
)

// Function to get the details of the CEO from the database.
func (db ElectionDatabase) GetCEO() (models.CEO, error) {
    ceoCollection := db.Session.DB(config.MongoDbName).C("ceo")
    ceo := models.CEO{}
    err := ceoCollection.Find(nil).One(&ceo)
    return ceo, err
}

// Function to insert a new CEO into the database.
func (db ElectionDatabase) InsertCEO(ceo *models.CEO) error {
    ceoCollection := db.Session.DB(config.MongoDbName).C("ceo")
    err := ceoCollection.Insert(&ceo)
    return err
}

// Function ot update the details of the CEO.
func (db ElectionDatabase) UpdateCEO(newCEO *models.CEO) error {
    ceoCollection := db.Session.DB(config.MongoDbName).C("ceo")
    err := ceoCollection.Update(bson.M{"username":"CEO"}, &newCEO)
    return err
}
