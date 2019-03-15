package utils

import (
    "crypto/tls"
    "fmt"
    "net/smtp"
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/models"
)

func GenerateMail(recipient *models.MailRecipient) string {
    mail := ""
    mail += fmt.Sprintf("From: %s\r\n", config.MailSenderEmailID)
    mail += fmt.Sprintf("To: %s\r\n", recipient.EmailID)
    mail += fmt.Sprintf("Subject: %s\r\n", config.MailSubject)
    mail += "\r\n"
    mail += fmt.Sprintf("Dear %s,\r\n", recipient.Name)
    mail += "\tThank you for registering as a voter for the upcoming Gymkhana Elections.\r\n"
    mail += "\tUse the following Authentication Code to complete your registration :\r\n"
    mail += fmt.Sprintf("\t\t%s\r\n",recipient.AuthCode)
    mail += "Regards,\r\nChief Election Officer,\r\nElection Commission,\r\nStudents' Senate,\r\nIIT Kanpur."
    return mail
}

func SendMailTo(recipient *models.MailRecipient) error {

    mail := GenerateMail(recipient)
    auth := smtp.PlainAuth("", config.MailSenderEmailID, config.MailSenderPassword, config.MailSMTPHost)
    tlsconfig := &tls.Config{
        InsecureSkipVerify:   true,
        ServerName:           config.MailSMTPHost,
    }

    connection, err := tls.Dial("tcp", config.MailSMTPHost+":"+config.MailSMTPPort, tlsconfig)
    if err != nil {
        return err
    }

    client, err := smtp.NewClient(connection, config.MailSMTPHost)
    if err != nil {
        return err
    }

    if err = client.Auth(auth); err != nil {
        return err
    }

    if err = client.Mail(config.MailSenderEmailID); err != nil {
        return err
    }

    if err = client.Rcpt(recipient.EmailID); err != nil {
        return err
    }

    writer, err := client.Data()
    if err != nil {
        return err
    }

    _, err = writer.Write([]byte(mail))
    if err != nil {
        return err
    }

    err = writer.Close()
    if err != nil {
        return err
    }

    client.Quit()

    return nil
}
