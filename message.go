package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"time"
)

// Email messages

type Address struct {
	Name, Email string
}

func (a *Address) String() string {
	if len(a.Name) == 0 {
		return a.Email
	} else {
		return Q_Encode(a.Name) + " <" + a.Email + ">"
	}
}

type Message struct {
	From    Address
	To      []Address
	Subject string
	Content MIMEPart
	CC      []Address
	BCC     []Address
	Headers map[string]string
}

func (m *Message) AllRecipients() []string {
	recipients := []string{}
	for _, address := range m.To {
		recipients = append(recipients, address.Email)
	}
	for _, address := range m.CC {
		recipients = append(recipients, address.Email)
	}
	for _, address := range m.BCC {
		recipients = append(recipients, address.Email)
	}
	return recipients
}

func (m *Message) Send(config SMTPConfig) error {
	err := smtp.SendMail(config.Server+":"+config.Port, config.GetSMTPAuth(), m.From.Email, m.AllRecipients(), []byte(m.String()))
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (m *Message) String() []byte {
	var buf bytes.Buffer
	buf.WriteString("From: ")
	buf.WriteString(m.From.String())
	buf.Write(crlf)
	for index, address := range m.To {
		if index == 0 {
			buf.WriteString("To: ")
		} else {
			buf.WriteString(",")
			buf.Write(crlf)
			buf.WriteString(" ")
		}
		buf.WriteString(address.String())
	}
	if len(m.To) != 0 {
		buf.Write(crlf)
	}
	for index, address := range m.CC {
		if index == 0 {
			buf.WriteString("Cc: ")
		} else {
			buf.WriteString(",")
			buf.Write(crlf)
			buf.WriteString(" ")
		}
		buf.WriteString(address.String())
	}
	if len(m.CC) != 0 {
		buf.Write(crlf)
	}
	for index, address := range m.BCC {
		if index == 0 {
			buf.WriteString("Bcc: ")
		} else {
			buf.WriteString(",")
			buf.Write(crlf)
			buf.WriteString(" ")
		}
		buf.WriteString(address.String())
	}
	if len(m.BCC) != 0 {
		buf.Write(crlf)
	}
	buf.WriteString("Subject: ")
	buf.WriteString(Q_Encode(m.Subject))
	buf.Write(crlf)
	buf.WriteString("Date: ")
	buf.WriteString(time.Now().Format("Mon, 2 Jan 2006 15:04:05 -0700"))
	buf.Write(crlf)
	for key, value := range m.Headers {
		buf.WriteString(key)
		buf.WriteString(": ")
		buf.WriteString(value)
		buf.Write(crlf)
	}
	buf.WriteString("MIME-Version: 1.0")
	buf.Write(crlf)
	Write(m.Content, &buf)
	return buf.Bytes()
}
