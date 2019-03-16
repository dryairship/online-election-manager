package models

import (
    "github.com/dryairship/online-election-manager/config"
)

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
    
    StudentSkeleton struct {
        Roll        string  `json:"roll"`
        Email       string  `json:"email"`
        Name        string  `json:"name"`
    }
)

func (voter Voter) GetMailRecipient() MailRecipient {
    return MailRecipient{
        Name:       voter.Name,
        EmailID:    voter.Email,
        AuthCode:   voter.AuthCode,
    }
}

func (skeleton StudentSkeleton) CreateVoter(authcode string) Voter {
    return Voter{
        Roll:       skeleton.Roll,
        Name:       skeleton.Name,
        Email:      skeleton.Email+config.MailSuffix,
        Password:   "",
        AuthCode:   authcode,
        BallotID:   "",
        Voted:      false,
    }
}
