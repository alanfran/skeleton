package main

import "fmt"

// Mailer sends emails.
type Mailer struct {
}

// NewMailer returns an initialized Mailer.
func NewMailer() *Mailer {

	return &Mailer{}
}

// Send sends an email to `addr`.
func (m Mailer) Send(addr, subj, body string) error {
	fmt.Println("[" + addr + "]  " + subj + ": " + body)
	return errNyi
}

// SendConfirmation sends a confirmation email.
func (m Mailer) SendConfirmation(addr, token string) error {
	m.Send(
		addr,
		"Welcome to "+appName+".",
		"Please activate your account.\nhttps://"+appURL+"/api/confirm?Token="+token,
	)

	return nil
}
