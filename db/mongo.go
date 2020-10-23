package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/dryairship/online-election-manager/config"
)

// Struct that is used to store the session of the database connection.
type ElectionDatabase struct {
	Session                     *mongo.Client
	VotersCollection            *mongo.Collection
	CandidatesCollection        *mongo.Collection
	CeoCollection               *mongo.Collection
	PostsCollection             *mongo.Collection
	BallotIdsCollection         *mongo.Collection
	SingleVoteResultsCollection *mongo.Collection
	ResultsCollection           *mongo.Collection
	StudentsCollection          *mongo.Collection
	VotesCollection             *mongo.Collection
}

// Function to establish database connection.
func ConnectToDatabase() (ElectionDatabase, error) {
	connectURL := fmt.Sprintf(
		"mongodb://%s:%s@%s/%s",
		url.QueryEscape(config.MongoUsername),
		url.QueryEscape(config.MongoPassword),
		config.MongoDialURL,
		config.MongoDbName,
	)

	var electionDb ElectionDatabase

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(connectURL))
	if err != nil {
		log.Printf("[ERROR] Cannot connect to Mongo. Error: %v\n", err)
		return electionDb, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("[ERROR] Cannot Ping Mongo. Error: %v\n", err)
		return electionDb, err
	} else {
		log.Println("[INFO] Successfully pinged Mongo")
		electionDb.Session = mongoClient
		electionDb.VotersCollection = mongoClient.Database(config.MongoDbName).Collection("voters")
		electionDb.CandidatesCollection = mongoClient.Database(config.MongoDbName).Collection("candidates")
		electionDb.CeoCollection = mongoClient.Database(config.MongoDbName).Collection("ceo")
		electionDb.PostsCollection = mongoClient.Database(config.MongoDbName).Collection("posts")
		electionDb.VotesCollection = mongoClient.Database(config.MongoDbName).Collection("votes")
		electionDb.StudentsCollection = mongoClient.Database(config.MongoDbName).Collection("students")
		electionDb.SingleVoteResultsCollection = mongoClient.Database(config.MongoDbName).Collection("singlevoteresults")
		electionDb.ResultsCollection = mongoClient.Database(config.MongoDbName).Collection("results")
		electionDb.BallotIdsCollection = mongoClient.Database(config.MongoDbName).Collection("ballotids")
	}

	return electionDb, nil
}

// Function to delete all entries from the database.
func (db ElectionDatabase) ResetDatabase() error {
	err := db.CandidatesCollection.Drop(context.Background())
	if err != nil {
		return err
	}

	err = db.VotesCollection.Drop(context.Background())
	if err != nil {
		return err
	}

	err = db.PostsCollection.Drop(context.Background())
	if err != nil {
		return err
	}

	err = db.SingleVoteResultsCollection.Drop(context.Background())
	if err != nil {
		return err
	}

	err = db.BallotIdsCollection.Drop(context.Background())
	if err != nil {
		return err
	}

	err = db.MarkAllVotersUnvoted()
	if err != nil {
		return err
	}

	// TODO: ceo, voters
	return nil
}
