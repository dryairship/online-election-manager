package models

type (
    User struct {
        Id          string  `json:"id"`
        Name        string  `json:"name"`
        Email       string  `json:"email"`
        Password    string  `json:"password"`
        Voted       bool    `json:"voted"`
    }
)