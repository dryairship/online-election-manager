package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dryairship/online-election-manager/models"
)

func (db ElectionDatabase) InsertResult(result *models.Result) error {
	_, err := db.ResultsCollection.InsertOne(context.Background(), result)
	return err
}

func (db ElectionDatabase) FindAllResults() ([]models.Result, error) {
	var results []models.Result

	cursor, err := db.ResultsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return results, err
	}

	err = cursor.All(context.Background(), &results)
	return results, err
}

func (db ElectionDatabase) ClearResults() error {
	err := db.ResultsCollection.Drop(context.Background())
	return err
}
