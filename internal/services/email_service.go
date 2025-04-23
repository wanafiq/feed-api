package services

import (
	"bytes"
	"github.com/wanafiq/feed-api/internal/config"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/email"
	"github.com/wanafiq/feed-api/internal/models"
	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
	"html/template"
)

type EmailService struct {
	config *config.Config
	logger *zap.SugaredLogger
}

type confirmationEmailVars struct {
	Username      string
	ActivationURL string
}

func NewEmailService(config *config.Config, logger *zap.SugaredLogger) *EmailService {
	return &EmailService{
		config: config,
		logger: logger,
	}
}

func (s *EmailService) Send(emailType string, data any, user *models.User) error {
	env := s.config.Env
	if env == constants.DevEnv {
		s.logger.Infow("skipping email sending", "env", env)
		return nil
	}

	var err error
	var tmpl *template.Template
	var subject string
	var htmlBody string

	switch emailType {
	case email.ConfirmationEmail:
		tmpl, err = s.getEmailTemplate(email.ConfirmationEmailTemplate)
		subject, err = s.getEmailSubject(tmpl)
		htmlBody, err = s.getEmailHtml(tmpl, data)
	default:
		return err
	}
	if err != nil {
		return err
	}

	if err := s.send(user.Email, subject, htmlBody); err != nil {
		return err
	}

	return nil
}

func (s *EmailService) send(to, subject, htmlBody string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.config.Smtp.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.AddAlternative("text/html", htmlBody)

	dialer := gomail.NewDialer(
		s.config.Smtp.Host,
		s.config.Smtp.Port,
		s.config.Smtp.Username,
		s.config.Smtp.Password,
	)
	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func (s *EmailService) getEmailTemplate(templateName string) (*template.Template, error) {
	tmpl, err := template.ParseFS(email.Templates, templateName)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s *EmailService) getEmailSubject(template *template.Template) (string, error) {
	subject := new(bytes.Buffer)

	err := template.ExecuteTemplate(subject, "subject", nil)
	if err != nil {
		return "", err
	}

	return subject.String(), nil
}

func (s *EmailService) getEmailHtml(template *template.Template, data any) (string, error) {
	body := new(bytes.Buffer)

	err := template.ExecuteTemplate(body, "body", data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}
