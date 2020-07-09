package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dryairship/online-election-manager/models"
)

// Function to find the candidate with that username in the database.
func (db ElectionDatabase) GetCandidate(username string) (models.Candidate, error) {
	candidate := models.Candidate{}
	err := db.CandidatesCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&candidate)
	return candidate, err
}

// Function to find all candidates in the database.
func (db ElectionDatabase) GetAllCandidates() ([]models.Candidate, error) {
	var candidates []models.Candidate

	cursor, err := db.CandidatesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return candidates, err
	}

	err = cursor.All(context.Background(), &candidates)
	return candidates, err
}

// Function to update the details of a candidate.
func (db ElectionDatabase) UpdateCandidate(username string, newCandidate *models.Candidate) error {
	_, err := db.CandidatesCollection.ReplaceOne(context.Background(), bson.M{"username": username}, newCandidate)
	return err
}

// Function to insert a new candidate in the database.
func (db ElectionDatabase) AddNewCandidate(candidate *models.Candidate) error {
	_, err := db.CandidatesCollection.InsertOne(context.Background(), candidate)
	return err
}

// Function to delete a candidate from the database.
func (db ElectionDatabase) EliminateCandidate(postId, roll string) error {
	username := fmt.Sprintf("P%sC%s", postId, roll)
	_, err := db.PostsCollection.UpdateOne(
		context.Background(),
		bson.M{"postid": postId},
		bson.M{
			"$pull": bson.M{
				"candidates": username,
			},
		},
	)
	return err
}
