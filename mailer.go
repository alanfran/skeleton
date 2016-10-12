package main

import (
  "fmt"
)

type Mailer struct {

}

func NewMailer() *Mailer {

  return &Mailer{}
}

func (m Mailer) Send(addr, subj, body string) error {
  return nyi
}

func (m Mailer) SendConfirmation(addr, token string) error {
  subject := `Welcome. Please activate your account.`
  fmt.Println("Sending mail to: " + addr)
  fmt.Println("Subject: " + subject)
  fmt.Println(token)

  return nil
}
