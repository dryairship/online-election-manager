package models

type SingleVoteCandidateResult struct {
	Name   string `json:"name"`
	Roll   string `json:"roll"`
	Count  int32  `json:"count"`
	Status string `json:"status"`
}

type SingleVoteResult struct {
	ID         string                      `json:"postId"`
	Name       string                      `json:"postName"`
	Candidates []SingleVoteCandidateResult `json:"candidates"`
}
