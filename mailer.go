package main

import (
	"fmt"
)

// Mailer sends emails.
type Mailer struct {
}

// NewMailer returns an initialized Mailer.
func NewMailer() *Mailer {

	return &Mailer{}
}

// Send sends an email to `addr`.
func (m Mailer) Send(addr, subj, body string) error {
	return errNyi
}

// SendConfirmation sends a confirmation email.
func (m Mailer) SendConfirmation(addr, token string) error {
	subject := `Welcome. Please activate your account.`
	fmt.Println("Sending mail to: " + addr)
	fmt.Println("Subject: " + subject)
	fmt.Println(token)

	return nil
}
