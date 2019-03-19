package models

type (
    Candidate struct {
        Roll        string  `json:"roll"`
        Name        string  `json:"name"`
        Email       string  `json:"email"`
        Username    string  `json:"username"`
        Password    string  `json:"password"`
        AuthCode    string  `json:"authcode"`
        PostID      string  `json:"postid"`
        Manifesto   string  `json:"manifesto"`
        PublicKey   string  `json:"publickey"`
        PrivateKey  string  `json:"privatekey"`
    }
    
    SimplifiedCandidate struct {
        Roll        string
        Name        string
        PublicKey   string
        Manifesto   string
    }
)

func (candidate Candidate) Simplify() SimplifiedCandidate {
    return SimplifiedCandidate {
        Roll:       candidate.Roll,
        Name:       candidate.Name,
        PublicKey:  candidate.PublicKey,
        Manifesto:  candidate.Manifesto,
    }
}
