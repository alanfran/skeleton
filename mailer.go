package main

import "fmt"

// Mailer sends emails.
type Mailer interface {
	Send(address, subject, body string) error
	SendConfirmation(address, token string) error
}

// MockMailer outputs mails to stdout.
type MockMailer struct {
	app *App
}

// NewMockMailer returns an initialized Mailer.
func NewMockMailer(a *App) *MockMailer {

	return &MockMailer{app: a}
}

// Send sends an email to `addr`.
func (m MockMailer) Send(addr, subj, body string) error {
	fmt.Println("[" + addr + "]  " + subj + ": " + body)
	return nil
}

// SendConfirmation sends a confirmation email.
func (m MockMailer) SendConfirmation(addr, token string) error {
	m.Send(
		addr,
		"Welcome to "+m.app.appName+".",
		"Please activate your account.\nhttps://"+m.app.appURL+"/api/confirm?Token="+token,
	)

	return nil
}
