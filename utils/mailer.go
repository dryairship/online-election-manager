package utils

import (
	"fmt"
	"net/smtp"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
)

var auth smtp.Auth

var verificationMailTemplate = `From: %s
To: %s
Subject: Verification Code for %s

Dear %s,
    Thank you for registering as %s for %s.

    Use the following verification code to complete your registration:
        %s

Regards,
%s
`

func AuthenticateMailer() {
	auth = smtp.PlainAuth("", config.MailSenderAuthID, config.MailSenderPassword, config.MailSMTPHost)
}

// Returns a string that represents the mail contents.
func GenerateMail(recipient *models.MailRecipient, role string) string {
	mail := fmt.Sprintf(verificationMailTemplate,
		config.MailSenderEmailID,
		recipient.EmailID,
		config.ElectionName,
		recipient.Name,
		role,
		config.ElectionName,
		recipient.AuthCode,
		config.MailSignature,
	)
	return mail
}

func SendMailTo(recipient *models.MailRecipient, role string) error {
	mail := []byte(GenerateMail(recipient, role))
	to := []string{recipient.EmailID}
	err := smtp.SendMail(config.MailSMTPHost+":"+config.MailSMTPPort, auth, config.MailSenderAuthID, to, mail)
	return err
}
