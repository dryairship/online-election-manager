package models

type (
	// Basic structure of a result as stored in the database.
	Result struct {
		PostID      string `json:"postid"`
		Candidate   string `json:"candidate"`
		Preference1 string `json:"preference1"`
		Preference2 string `json:"preference2"`
		Preference3 string `json:"preference3"`
	}
)
