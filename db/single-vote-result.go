package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dryairship/online-election-manager/models"
)

// Function to insert a new result into the database.
func (db ElectionDatabase) InsertSingleVoteResult(result *models.SingleVoteResult) error {
	_, err := db.SingleVoteResultsCollection.InsertOne(context.Background(), result)
	return err
}

func (db ElectionDatabase) FindAllSingleVoteResults() ([]models.SingleVoteResult, error) {
	var results []models.SingleVoteResult

	cursor, err := db.SingleVoteResultsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return results, err
	}

	err = cursor.All(context.Background(), &results)
	return results, err
}

func (db ElectionDatabase) ClearSingleVoteResults() error {
	err := db.SingleVoteResultsCollection.Drop(context.Background())
	return err
}
