package utils

import (
    "crypto/tls"
    "fmt"
    "net/smtp"
    "github.com/dryairship/online-election-manager/config"
    "github.com/dryairship/online-election-manager/models"
)

// Returns a string that represents the mail contents.
func GenerateMail(recipient *models.MailRecipient) string {
    mail := fmt.Sprintf("From: %s\r\n", config.MailSenderEmailID) +
            fmt.Sprintf("To: %s\r\n", recipient.EmailID) +
            fmt.Sprintf("Subject: %s\r\n", config.MailSubject) +
            "\r\n" +
            fmt.Sprintf("Dear %s,\r\n", recipient.Name) +
            "\tThank you for registering as a voter for the upcoming Gymkhana Elections.\r\n" +
            "\tUse the following Authentication Code to complete your registration :\r\n" +
            fmt.Sprintf("\t\t%s\r\n",recipient.AuthCode) +
            "Regards,\r\nChief Election Officer,\r\nElection Commission,\r\nStudents' Senate,\r\nIIT Kanpur."
    return mail
}

func SendMailTo(recipient *models.MailRecipient) error {
    // Generate the mail contents for this recipient.
    mail := GenerateMail(recipient)
    
    // Set up authentication details.
    auth := smtp.PlainAuth("", config.MailSenderEmailID, config.MailSenderPassword, config.MailSMTPHost)
    tlsconfig := &tls.Config{
        InsecureSkipVerify:   true,
        ServerName:           config.MailSMTPHost,
    }
    
    // Dial up the SMTP host for a connection.
    connection, err := tls.Dial("tcp", config.MailSMTPHost+":"+config.MailSMTPPort, tlsconfig)
    if err != nil {
        return err
    }
    
    // Get the client object from the connection.
    client, err := smtp.NewClient(connection, config.MailSMTPHost)
    if err != nil {
        return err
    }
    
    // Authenticate the client.
    if err = client.Auth(auth); err != nil {
        return err
    }
    
    // Set up the Mail ID of the sender.
    if err = client.Mail(config.MailSenderEmailID); err != nil {
        return err
    }
    
    // Set up the Mail ID of the recipient.
    if err = client.Rcpt(recipient.EmailID); err != nil {
        return err
    }
    
    // Get the stream to write the mail contents.
    writer, err := client.Data()
    if err != nil {
        return err
    }
    
    // Wtite the mail contents to the stream.
    _, err = writer.Write([]byte(mail))
    if err != nil {
        return err
    }
    
    // Close the stream.
    err = writer.Close()
    if err != nil {
        return err
    }
    
    // Close the client connection.
    client.Quit()
    
    return nil
}
