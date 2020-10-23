package models

import "fmt"

type (
	// Basic structure of a vote as stored in the database.
	Vote struct {
		PostID string `json:"postid"`
		Data   string `json:"data"`
	}

	// Basic structure of a parsed vote as stored in the database.
	ParsedVote struct {
		PostID      string `json:"postid"`
		BallotID    string `json:"ballotid"`
		Preference1 string `json:"preference1"`
		Preference2 string `json:"preference2"`
		Preference3 string `json:"preference3"`
	}

	// Struct to represent the received Ballot ID.
	BallotID struct {
		PostID       string
		BallotString string
	}

	// Struct to represent the received vote from the user.
	ReceivedVote struct {
		PostID       int    `json:"PostID"`
		BallotString string `json:"BallotString"`
		VoteData     string `json:"VoteData"`
	}

	UsedSingleVoteBallotID struct {
		BallotString string
		Roll         string
		Name         string
	}

	UsedBallotID struct {
		BallotString string `json:"ballotString"`
		Preference1  string `json:"preference1"`
		Preference2  string `json:"preference2"`
		Preference3  string `json:"preference3"`
	}
)

// Function to get the actual data of the vote from it.
func (receivedVote ReceivedVote) GetVote() Vote {
	return Vote{
		PostID: fmt.Sprintf("%d", receivedVote.PostID),
		Data:   receivedVote.VoteData,
	}
}

// Function to get the ballot ID fron a vote.
func (receivedVote ReceivedVote) GetBallotID() BallotID {
	return BallotID{
		PostID:       fmt.Sprintf("%d", receivedVote.PostID),
		BallotString: receivedVote.BallotString,
	}
}
