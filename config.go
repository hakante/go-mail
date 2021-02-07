package mail

import (
	"net/smtp"
)

type SMTPConfig struct {
	Username, Password, Server, Port string
}

func (config SMTPConfig) GetSMTPAuth() smtp.Auth {
	return smtp.PlainAuth("", config.Username, config.Password, config.Server)
}
