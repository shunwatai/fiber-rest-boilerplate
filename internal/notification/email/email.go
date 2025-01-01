package email

import (
	"golang-api-starter/internal/config"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"html/template"

	m "net/mail"

	"github.com/wneessen/go-mail"
)

var cfg = config.Cfg

type IQueue interface {
	QueueMsg() error
}

type EmailInfo struct {
	From          string
	To            []string
	Cc            []string
	Bcc           []string
	MsgMeta       map[string]interface{} // includes Subject & other data if use Template
	MsgContent    string                 // required for SimpleEmail
	TmplFilePaths []string               // required for TemplateEmail
}

// ref: https://stackoverflow.com/a/66624104
func validateEmail(email string) bool {
	_, err := m.ParseAddress(email)
	return err == nil
}

func (e *EmailInfo) checkInfo() error {
	if len(e.To) == 0 {
		return logger.Errorf("Recipient is required")
	}

	for _, recipientEmail := range e.To {
		if !validateEmail(recipientEmail) {
			return logger.Errorf("%s is not a valid email...", recipientEmail)
		}
	}

	// if send email by template, no need to check for "SimpleEmail"'s e.Subject e.Message
	if len(e.TmplFilePaths) > 0 {
		return nil
	}

	if len(e.MsgMeta["subject"].(string)) == 0 {
		return logger.Errorf("Subject is required")
	}
	if len(e.MsgContent) == 0 {
		return logger.Errorf("Cannot send empty message")
	}

	return nil
}

func (e *EmailInfo) getClient() (*mail.Client, error) {
	client, err := mail.NewClient(
		cfg.Notification.Smtp.Host,
		mail.WithPort(cfg.Notification.Smtp.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.Notification.Smtp.User),
		mail.WithPassword(cfg.Notification.Smtp.Pass),
	)

	if err != nil {
		return nil, logger.Errorf("failed to create mail client: %s", err)
	}

	return client, nil
}

func (e *EmailInfo) initMessageWithInfo() (*mail.Msg, error) {
	m := mail.NewMsg()

	e.From = cfg.Notification.Smtp.User
	if err := m.From(e.From); err != nil {
		return nil, logger.Errorf("failed to set From address: %s", err)
	}
	if err := m.To(e.To...); err != nil {
		return nil, logger.Errorf("failed to set To address: %s", err)
	}
	m.Subject(e.MsgMeta["subject"].(string))

	return m, nil
}

func SimpleEmail(info EmailInfo) error {
	var (
		mailClient *mail.Client
		m          *mail.Msg
		err        error
	)

	if err := info.checkInfo(); err != nil {
		return err
	}

	if m, err = info.initMessageWithInfo(); err != nil {
		return err
	}
	m.SetBodyString(mail.TypeTextHTML, info.MsgContent)

	if mailClient, err = info.getClient(); err != nil {
		return err
	}

	if err := mailClient.DialAndSend(m); err != nil {
		return logger.Errorf("failed to send mail: %s", err)
	}

	return nil
}

func TemplateEmail(info EmailInfo) error {
	if info.TmplFilePaths == nil {
		return logger.Errorf("Cannot send email without template, please assign the Template to info")
	}
	tmpl := template.Must(template.ParseFiles(info.TmplFilePaths...))

	var (
		mailClient *mail.Client
		m          *mail.Msg
		err        error
	)

	if err := info.checkInfo(); err != nil {
		return err
	}

	if m, err = info.initMessageWithInfo(); err != nil {
		return err
	}

	if err := m.SetBodyHTMLTemplate(tmpl, info.MsgMeta); err != nil {
		return err
	}

	if mailClient, err = info.getClient(); err != nil {
		return err
	}

	if err := mailClient.DialAndSend(m); err != nil {
		return logger.Errorf("failed to send mail: %s", err)
	}

	return nil
}
