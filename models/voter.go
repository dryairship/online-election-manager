package models

type (
    Voter struct {
        Roll        string  `json:"roll"`
        Name        string  `json:"name"`
        Email       string  `json:"email"`
        Password    string  `json:"password"`
        AuthCode    string  `json:"authcode"`
        BallotID    string  `json:"ballotid"`
        Voted       bool    `json:"voted"`
    }

    MailRecipient struct {
        Name        string
        EmailID     string
        AuthCode    string
    }
) 
