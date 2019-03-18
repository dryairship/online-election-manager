package models

type (
    Post struct {
        PostID          string      `json:"postid"`
        PostName        string      `json:"postname"`
        VoterRegex      string      `json:"regex"`
        Candidates      []string    `json:"candidates"`
    }
    
    VotablePost struct {
        PostID          string
        PostName        string
        Candidates      []string
    }
)

func (post Post) ConvertToVotablePost() VotablePost {
    return VotablePost {
        PostID:         post.PostID,
        PostName:       post.PostName,
        Candidates:     post.Candidates,
    }
}
