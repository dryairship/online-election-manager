package db

import (
	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
)

// Function to insert a new vote into the database.
func (db ElectionDatabase) InsertVote(vote *models.Vote) error {
	votesCollection := db.Session.DB(config.MongoDbName).C("votes")
	err := votesCollection.Insert(&vote)
	return err
}

// Function to insert a new parsed vote into the database.
func (db ElectionDatabase) InsertParsedVote(parsedVote *models.ParsedVote) error {
	parsedVotesCollection := db.Session.DB(config.MongoDbName).C("parsedvotes")
	err := parsedVotesCollection.Insert(&parsedVote)
	return err
}

// Function to get all the votes from the database.
func (db ElectionDatabase) GetVotes() ([]models.Vote, error) {
	votesCollection := db.Session.DB(config.MongoDbName).C("votes")
	var votes []models.Vote
	err := votesCollection.Find(nil).All(&votes)
	return votes, err
}
