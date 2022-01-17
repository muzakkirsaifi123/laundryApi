package smtp

import (
	"fmt"
	"net/smtp"

	"github.com/aymerick/raymond"
	config "github.com/spf13/viper"
)

func init() {
	raymond.RegisterHelper(
		"resetButton", func(uid int, token string) raymond.SafeString {
			link := fmt.Sprintf("http://localhost:9090/passwordReset/%d/%s", uid, token)
			return raymond.SafeString(
				fmt.Sprintf(
					"<a href=%s target=\"_blank\">Reset Password</a>", raymond.Escape(link),
				),
			)
		},
	)
}

type Smtp struct {
	host                  string
	port                  int
	username              string
	password              string
	passwordResetTemplate *raymond.Template
}

type ResetPasswordPayload struct {
	UID       int
	Token     string
	FirstName string
	Email     string
}

func (s Smtp) SendResetEmail(payload ResetPasswordPayload) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subj := "Subject: Reset Password!\n"
	page := s.passwordResetTemplate.MustExec(payload)
	msg := fmt.Sprintf("%s%s\n%s", subj, mime, page)
	smtpHostPort := fmt.Sprintf("%s:%d", s.host, s.port)
	return smtp.SendMail(
		smtpHostPort, smtp.PlainAuth("", s.username, s.password, s.host), s.username, []string{payload.Email},
		[]byte(msg),
	)
}

func New() (Smtp, error) {
	fmt.Printf("creating smtp server")
	tpl, err := raymond.ParseFile("./templates/passwordResetEmail.html")
	return Smtp{
		host:                  config.GetString("smtp_host"),
		port:                  config.GetInt("smtp_port"),
		username:              config.GetString("smtp_username"),
		password:              config.GetString("smtp_password"),
		passwordResetTemplate: tpl,
	}, err
}
