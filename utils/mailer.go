package utils

import (
	"fmt"
	"net/smtp"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/models"
)

var auth smtp.Auth

func AuthenticateMailer() {
	auth = smtp.PlainAuth("", config.MailSenderUsername, config.MailSenderPassword, config.MailSMTPHost)
}

// Returns a string that represents the mail contents.
func GenerateMail(recipient *models.MailRecipient, role string) string {
	mail := fmt.Sprintf("From: %s\r\n", config.MailSenderUsername+"@iitk.ac.in") +
		fmt.Sprintf("To: %s\r\n", recipient.EmailID) +
		"Subject: Mid-Term and By-Elections Verification Code\r\n" +
		"\r\n" +
		fmt.Sprintf("Dear %s,\r\n", recipient.Name) +
		fmt.Sprintf("\tThank you for registering as %s for the Gymkhana Elections.\r\n", role) +
		"\tUse the following Authentication Code to complete your registration :\r\n" +
		fmt.Sprintf("\t\t%s\r\n", recipient.AuthCode) +
		"Regards,\r\nChief Election Officer,\r\nElection Commission,\r\nIIT Kanpur.\r\n"
	return mail
}

func SendMailTo(recipient *models.MailRecipient, role string) error {
	mail := []byte(GenerateMail(recipient, role))
	to := []string{recipient.EmailID}
	return smtp.SendMail(config.MailSMTPHost+":"+config.MailSMTPPort, auth, config.MailSenderUsername, to, mail)
}
