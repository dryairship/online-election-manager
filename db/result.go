package db

import (
	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
)

// Function to insert a new result into the database.
func (db ElectionDatabase) InsertResult(result *models.Result) error {
	resultsCollection := db.Session.DB(config.MongoDbName).C("results")
	err := resultsCollection.Insert(&result)
	return err
}
