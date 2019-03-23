package models

type (
    Vote struct {
        PostID          string  `json:"postid"`
        Data            string  `json:"data"`
    }
    
    BallotID struct {
        PostID          string
        BallotString    string
    }
    
    ReceivedVote struct {
        PostID          int     `json:"PostID"`
        BallotString    string  `json:"BallotString"`
        VoteData        string  `json:"VoteData"`
    }
)

func (receivedVote ReceivedVote) GetVote() Vote {
    return Vote {
        PostID: string(receivedVote.PostID),
        Data:   receivedVote.VoteData,
    }
}

func (receivedVote ReceivedVote) GetBallotID() BallotID {
    return BallotID {
        PostID:         string(receivedVote.PostID),
        BallotString:   receivedVote.BallotString,
    }
}

