package db

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dryairship/online-election-manager/models"
)

// Function to insert a new vote into the database.
func (db ElectionDatabase) InsertVote(vote *models.Vote) error {
	_, err := db.VotesCollection.InsertOne(context.Background(), vote)
	return err
}

// Function to insert a new parsed vote into the database.
// TODO
func (db ElectionDatabase) InsertParsedVote(parsedVote *models.ParsedVote) error {
	//parsedVotesCollection := db.Session.DB(config.MongoDbName).C("parsedvotes")
	//err := parsedVotesCollection.Insert(&parsedVote)
	return nil
}

// Function to get all the votes from the database.
func (db ElectionDatabase) GetVotes() ([]models.Vote, error) {
	var votes []models.Vote

	cursor, err := db.VotesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return votes, err
	}

	err = cursor.All(context.Background(), &votes)
	return votes, err
}

func (db ElectionDatabase) DeleteAllVotes() error {
	err := db.VotesCollection.Drop(context.Background())
	return err
}

func ToInterfaceArray(s interface{}) []interface{} {
	v := reflect.ValueOf(s)
	intf := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		intf[i] = v.Index(i).Interface()
	}
	return intf
}
func (db ElectionDatabase) AddBallotIDs(ballotIds []models.UsedBallotID) error {
	_, err := db.BallotIdsCollection.InsertMany(context.Background(), ToInterfaceArray(ballotIds))
	return err
}
