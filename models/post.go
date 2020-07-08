package models

type (
	// Basic structure of the posts as stored in the database.
	Post struct {
		PostID     string   `json:"postid"`
		PostName   string   `json:"postname"`
		VoterRegex string   `json:"regex"`
		Candidates []string `json:"candidates"`
		Resolved   bool     `json:"resolved"`
	}

	// Structure of posts returned by the appropriate API call.
	VotablePost struct {
		PostID     string
		PostName   string
		Candidates []string
	}
)

// Function to strip out regex from the data of the post.
func (post Post) ConvertToVotablePost() VotablePost {
	return VotablePost{
		PostID:     post.PostID,
		PostName:   post.PostName,
		Candidates: post.Candidates,
	}
}
