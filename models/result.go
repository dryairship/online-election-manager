package models

import "strconv"

type (
	// Basic structure of a result as stored in the database.
	Result struct {
		PostID      string `json:"postid"`
		Candidate   string `json:"candidate"`
		Preference1 string `json:"preference1"`
		Preference2 string `json:"preference2"`
		Preference3 string `json:"preference3"`
	}

	NumericResult struct {
		PostID      int
		Candidate   string
		Preference1 int
		Preference2 int
		Preference3 int
	}

	ResultSorter []NumericResult
)

func (result *Result) ToNumericResult() NumericResult {
	postId, _ := strconv.Atoi(result.PostID)
	p1, _ := strconv.Atoi(result.Preference1)
	p2, _ := strconv.Atoi(result.Preference2)
	p3, _ := strconv.Atoi(result.Preference3)
	return NumericResult{
		PostID:      postId,
		Candidate:   result.Candidate,
		Preference1: p1,
		Preference2: p2,
		Preference3: p3,
	}
}

func (a ResultSorter) Len() int {
	return len(a)
}

func (a ResultSorter) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ResultSorter) Less(i, j int) bool {
	if a[i].PostID != a[j].PostID {
		return a[i].PostID < a[j].PostID
	}

	if a[i].Preference1 > a[j].Preference1 {
		return true
	} else if a[i].Preference1 == a[j].Preference1 {
		if a[i].Preference2 > a[j].Preference2 {
			return true
		} else if a[i].Preference2 == a[j].Preference2 {
			return a[i].Preference3 > a[j].Preference3
		}
	}
	return false
}
