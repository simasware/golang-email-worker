package services

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/spf13/viper"
	"simasware.com.br/email-worker/models"
)

func setupHeaders(mailRequest models.SendEmailRequest) map[string]string {
	headers := make(map[string]string)
	headers["From"] = viper.GetString("MAIL.MAIL_FROM")
	headers["To"] = mailRequest.Recipient
	headers["Subject"] = mailRequest.Subject
	headers["Mime"] = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	return headers
}

func setUpBody(mailRequest models.SendEmailRequest) string {
	headers := setupHeaders(mailRequest)
	messageBody := ""

	for k, v := range headers {
		messageBody += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	messageBody += "\r\n" + mailRequest.Body

	return messageBody
}

func SecureMail(mailRequest models.SendEmailRequest) (bool, error) {
	from := mail.Address{"", viper.GetString("MAIL.MAIL_FROM")}
	to := mail.Address{"", mailRequest.Recipient}

	message := setUpBody(mailRequest)
	servername := viper.GetString("MAIL.SMTP") + ":" + viper.GetString("MAIL.SMTP_PORT")
	auth := smtp.PlainAuth("", viper.GetString("MAIL.MAIL_FROM"), viper.GetString("MAIL.MAIL_PASSWORD"), viper.GetString("MAIL.SMTP"))

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         viper.GetString("MAIL.SMTP"),
	}

	connection, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return false, err
	}

	client, err := smtp.NewClient(connection, viper.GetString("MAIL.SMTP"))
	if err != nil {
		return false, err
	}

	if err = client.Auth(auth); err != nil {
		return false, err
	}

	if err = client.Mail(from.Address); err != nil {
		return false, err
	}

	if err = client.Rcpt(to.Address); err != nil {
		return false, err
	}

	w, err := client.Data()
	if err != nil {
		return false, err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return false, err
	}

	err = w.Close()
	if err != nil {
		return false, err
	}

	client.Quit()

	return true, nil
}
