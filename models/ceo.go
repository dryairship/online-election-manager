package models

import "github.com/dryairship/online-election-manager/config"

type (
    CEO struct {
        Roll        string  `json:"roll"`
        Name        string  `json:"name"`
        Email       string  `json:"email"`
        Username    string  `json:"username"`
        Password    string  `json:"password"`
        AuthCode    string  `json:"authcode"`
        PublicKey   string  `json:"publickey"`
        PrivateKey  string  `json:"privatekey"`
    }
)

func (skeleton StudentSkeleton) CreateCEO(authcode string) CEO {
    return CEO {
        Roll:       skeleton.Roll,
        Name:       skeleton.Name,
        Email:      skeleton.Email+config.MailSuffix,
        Username:   "CEO",
        Password:   "",
        AuthCode:   authcode,
        PublicKey:  "",
        PrivateKey: "",
    }
}

func (ceo CEO) GetMailRecipient() MailRecipient {
    return MailRecipient {
        Name:       ceo.Name,
        EmailID:    ceo.Email,
        AuthCode:   ceo.AuthCode,
    }
}
