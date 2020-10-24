package models

type (
	// Basic structure of a result as stored in the database.
	CandidateResult struct {
		Name        string `json:"name"`
		Roll        string `json:"roll"`
		Preference1 int32  `json:"preference1"`
		Preference2 int32  `json:"preference2"`
		Preference3 int32  `json:"preference3"`
		Status      string `json:"status"`
	}

	Result struct {
		ID         string            `json:"postId"`
		Name       string            `json:"postName"`
		Candidates []CandidateResult `json:"candidates"`
	}
)
