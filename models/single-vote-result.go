package models

type CandidateResult struct {
	Name   string `json:"name"`
	Roll   string `json:"roll"`
	Count  int32  `json:"count"`
	Status string `json:"status"`
}

type SingleVoteResult struct {
	ID         string            `json:"postId"`
	Name       string            `json:"postName"`
	Candidates []CandidateResult `json:"candidates"`
}
