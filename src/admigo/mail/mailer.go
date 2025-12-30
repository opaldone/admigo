package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"admigo/applog"
	"admigo/common"
	"admigo/config"
)

type EmMessage struct {
	From    string
	Subject string
	To      string
	Email   string
	Key     string
	Link    string
}

// SendEmailUserRegistration sends email
func SendEmailUserRegistration(email string) {
	e := config.Env(false)

	message := &EmMessage{
		From:    e.Mail.From,
		To:      email,
		Subject: "User Registration on Admigo",
		Email:   email,
		Key:     common.Hapas(email),
		Link:    fmt.Sprintf("%s:%d", e.Mail.GotoURL, e.Port),
	}

	var body bytes.Buffer

	GenerateMail(&body, message, "lays/layout_mail", "mail/reguser")

	auth := smtp.PlainAuth("",
		e.Mail.Username,
		e.Mail.Password,
		e.Mail.Host,
	)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         e.Mail.Host,
	}

	servername := fmt.Sprintf("%s:%d", e.Mail.Host, e.Mail.Port)

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		applog.Danger("mailer.go tls.Dial", err)
		return
	}

	c, err := smtp.NewClient(conn, e.Mail.Host)
	if err != nil {
		applog.Danger("mailer.go smtp.NewClient", err)
		return
	}

	if err = c.Auth(auth); err != nil {
		applog.Danger("mailer.go c.Auth", err)
		return
	}

	// To && From
	if err = c.Mail(message.From); err != nil {
		applog.Danger("mailer.go c.Mail", err)
		return
	}

	if err = c.Rcpt(message.To); err != nil {
		applog.Danger("mailer.go c.Rcpt", err)
		return
	}

	// Data
	w, err := c.Data()
	if err != nil {
		applog.Danger("mailer.go c.Data()", err)
		return
	}

	_, err = w.Write(body.Bytes())
	if err != nil {
		applog.Danger("mailer.go w.Write", err)
		return
	}

	err = w.Close()
	if err != nil {
		applog.Danger("mailer.go w.Close", err)
	}

	c.Quit()
}
