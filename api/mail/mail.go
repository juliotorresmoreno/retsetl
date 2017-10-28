package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

//Send Send an email using smtp protocol, only works with gmail accounts
//To work you must enable the sending of messages from less secure applications
//https://myaccount.google.com/u/1/lesssecureapps?rfn=27&rfnc=1&eid=6335218851849783228&et=0&asae=2&pli=1
//You must also enable the sending of imap messages
//
//@param subject string Message Subject
//@param to      string Message recipient
//@param from    string Message sender
//@param pass    string Sending user password
//@param body    string message body
func Send(subject, to, from, pass, body string) error {
	fromAddress := mail.Address{Address: from}
	toAddress := mail.Address{Address: to}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = fromAddress.String()
	headers["To"] = toAddress.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := "smtp.gmail.com:465"

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return nil
}
