package mailgunx

import "github.com/mailgun/mailgun-go/v4"

type MailgunOpt struct {
	Domain  string `help:"kafka address" required:"" env:"DOMAIN"`
	APIKey  string `help:"kafka address" required:"" env:"APIKEY"`
	BaseURL string `help:"kafka address" required:"" env:"BASE_URL"`
}

func NewMailgun(opt MailgunOpt) *mailgun.MailgunImpl {
	return mailgun.NewMailgun(opt.Domain, opt.APIKey)
}
