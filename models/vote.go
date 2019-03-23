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
)

