package models

import "fmt"

type (
    // Basic structure of a vote as stored in the databse.
    Vote struct {
        PostID          string  `json:"postid"`
        Data            string  `json:"data"`
    }
    
    // Struct to represent the received Ballot ID.
    BallotID struct {
        PostID          string
        BallotString    string
    }
    
    // Struct to represent the received vote from the user.
    ReceivedVote struct {
        PostID          int     `json:"PostID"`
        BallotString    string  `json:"BallotString"`
        VoteData        string  `json:"VoteData"`
    }
)

// Function to get the actual data of the vote from it.
func (receivedVote ReceivedVote) GetVote() Vote {
    return Vote {
        PostID: fmt.Sprintf("%d",receivedVote.PostID),
        Data:   receivedVote.VoteData,
    }
}

// Function to get the ballot ID fron a vote.
func (receivedVote ReceivedVote) GetBallotID() BallotID {
    return BallotID {
        PostID:         fmt.Sprintf("%d",receivedVote.PostID),
        BallotString:   receivedVote.BallotString,
    }
}

