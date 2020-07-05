package db

import (
	"github.com/dryairship/online-election-manager/models"
)

// Function to insert a new result into the database.
// TODO
func (db ElectionDatabase) InsertResult(result *models.Result) error {
	//resultsCollection := db.Session.DB(config.MongoDbName).C("results")
	//err := resultsCollection.Insert(&result)
	return nil
}

// TODO
func (db ElectionDatabase) FindAllResults() ([]models.Result, error) {
	//resultsCollection := db.Session.DB(config.MongoDbName).C("results")
	var results []models.Result
	//err := resultsCollection.Find(nil).All(&results)
	return results, nil
}

// TODO
func (db ElectionDatabase) ClearResults() error {
	//resultsCollection := db.Session.DB(config.MongoDbName).C("results")
	//_, err := resultsCollection.RemoveAll(nil)
	return nil
}
