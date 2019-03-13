package models

type (
    Candidate struct {
        Roll        string  `json:"roll"`
        Name        string  `json:"name"`
        Email       string  `json:"email"`
        Password    string  `json:"password"`
        AuthCode    string  `json:"authcode"`
        Manifesto   string  `json:"manifesto"`
        PublicKey   string  `json:"publickey"`
        PrivateKey  string  `json:"privatekey"`
    }
) 