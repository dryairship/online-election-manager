package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dryairship/online-election-manager/models"
)

// Function to get the details of the CEO from the database.
func (db ElectionDatabase) GetCEO() (models.CEO, error) {
	ceo := models.CEO{}
	err := db.CeoCollection.FindOne(context.Background(), bson.M{}).Decode(&ceo)
	return ceo, err
}

// Function to insert a new CEO into the database.
func (db ElectionDatabase) InsertCEO(ceo *models.CEO) error {
	_, err := db.CeoCollection.InsertOne(context.Background(), &ceo)
	return err
}

// Function ot update the details of the CEO.
func (db ElectionDatabase) UpdateCEO(newCEO *models.CEO) error {
	_, err := db.CeoCollection.ReplaceOne(context.Background(), bson.M{"username": "CEO"}, &newCEO)
	return err
}
