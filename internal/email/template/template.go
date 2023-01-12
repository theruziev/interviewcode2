package template

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed welcome.html
var welcomeEmailHTML string

//go:embed reset-password.html
var resetPasswordHTML string

var welcomeEmailTemplate = template.Must(template.New("welcome").Parse(welcomeEmailHTML))
var resetPasswordTemplate = template.Must(template.New("reset-password").Parse(resetPasswordHTML))

func WelcomeEmail(name, activationLink string) (string, error) {
	data := struct {
		Name           string
		ActivationLink string
	}{
		Name:           name,
		ActivationLink: activationLink,
	}

	strBuffer := bytes.NewBufferString("")

	if err := welcomeEmailTemplate.Execute(strBuffer, data); err != nil {
		return "", err
	}

	return strBuffer.String(), nil
}

func ResetPasswordEmail(name, resetLink string) (string, error) {
	data := struct {
		Name      string
		ResetLink string
	}{
		Name:      name,
		ResetLink: resetLink,
	}

	strBuffer := bytes.NewBufferString("")
	if err := resetPasswordTemplate.Execute(strBuffer, data); err != nil {
		return "", err
	}

	return strBuffer.String(), nil
}
