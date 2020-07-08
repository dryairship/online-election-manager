package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dryairship/online-election-manager/models"
)

// Function to find the basic details of a student from its roll number.
func (db ElectionDatabase) FindStudentSkeleton(roll string) (models.StudentSkeleton, error) {
	skeleton := models.StudentSkeleton{}
	err := db.StudentsCollection.FindOne(context.Background(), bson.M{"roll": roll}).Decode(&skeleton)
	return skeleton, err
}

// Function to find a voter in the database.
func (db ElectionDatabase) FindVoter(roll string) (models.Voter, error) {
	voter := models.Voter{}
	err := db.VotersCollection.FindOne(context.Background(), bson.M{"roll": roll}).Decode(&voter)
	return voter, err
}

// Function to find voters who have voted.
func (db ElectionDatabase) FindVotedVoters() ([]models.Voter, error) {
	var voters []models.Voter
	cursor, err := db.VotersCollection.Find(context.Background(), bson.M{"voted": true})
	if err != nil {
		return voters, err
	}

	err = cursor.All(context.Background(), &voters)
	return voters, err
}

// Function to insert a new voter into the database.
func (db ElectionDatabase) AddNewVoter(voter *models.Voter) error {
	_, err := db.VotersCollection.InsertOne(context.Background(), voter)
	return err
}

// Function to update the details of a voter.
func (db ElectionDatabase) UpdateVoter(roll string, newVoter *models.Voter) error {
	_, err := db.VotersCollection.ReplaceOne(context.Background(), bson.M{"roll": roll}, newVoter)
	return err
}

func (db ElectionDatabase) MarkAllVotersUnvoted() error {
	emptyArray := make([]models.BallotID, 0)
	_, err := db.VotersCollection.UpdateMany(
		context.Background(),
		bson.M{},
		bson.M{"$set": bson.M{
			"voted":    false,
			"ballotid": emptyArray,
		}},
	)
	return err
}
