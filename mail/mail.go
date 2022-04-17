package mail

import (
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"
)

// Mailer exported Mailer struct
type Mailer struct {
	senderName  string
	senderEmail string
	dialer      IDialer
}

// NewMailer construct the mailer object
func NewMailer(
	sendername string, senderemail string, host string, port int,
	username string, password string) *Mailer {
	dialer := gomail.NewDialer(host, port, username, password)
	return &Mailer{
		senderName: sendername, dialer: dialer,
		senderEmail: senderemail}
}

// Send sends all the SendableMessages using a single connection
func (mailer *Mailer) Send(msgs ...SendableMessage) error {
	messages := mailer.buildMessages(msgs...)
	if err := mailer.dialer.DialAndSend(messages...); err != nil {
		return err
	}
	return nil
}

func (mailer *Mailer) buildMessages(msgs ...SendableMessage) []*gomail.Message {
	messages := []*gomail.Message{}
	for _, msg := range msgs {
		message := mailer.buildMessage(msg)
		messages = append(messages, message)
	}
	return messages
}

func (mailer *Mailer) buildMessage(msg SendableMessage) *gomail.Message {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", mailer.senderEmail, mailer.senderName)
	to := msg.GetRecipients()
	setRecipients(m, "To", to...)
	cc := msg.GetCc()
	setRecipients(m, "Cc", cc...)
	bcc := msg.GetBCc()
	setRecipients(m, "BCc", bcc...)
	subject := msg.GetSubject()
	m.SetHeader("Subject", subject)
	setMessageBody(m, msg)
	setAttachments(m, msg)
	setHeaders(m, msg)
	return m
}

func setRecipients(m *gomail.Message, name string, recipients ...string) {
	for _, rcpt := range recipients {
		m.SetAddressHeader(name, rcpt, "")
	}
}

func setMessageBody(m *gomail.Message, msg SendableMessage) {
	textMessage := msg.GetTextMessage()
	plainTextType := "text/plain"

	htmlMessage := msg.GetHTMLMessage()
	htmlType := "text/html"

	// email clients receive emails with alternatives in a way that the last
	// part is given priority, I want to avoid a condition where an empty email
	// message is rendered over a plain text message that isn't empty
	if textMessage == "" {
		m.SetBody(htmlType, htmlMessage)
	} else if htmlMessage != "" {
		m.SetBody(plainTextType, textMessage)
		m.AddAlternative(htmlType, htmlMessage)
	} else {
		m.SetBody(plainTextType, textMessage)
	}
}

func setAttachments(m *gomail.Message, msg SendableMessage) {
	for _, attachment := range msg.GetAttachments() {
		m.Attach(attachment)
	}
}

func setHeaders(m *gomail.Message, msg SendableMessage) {
	headers := msg.GetHeaders()
	for key, value := range headers {
		m.SetHeader(key, value...)
	}
}

// Message to create an instance of this struct please use InitializeMessage
// function use it as mails.InitializeMessage()
type Message struct {
	Subject     string
	HTMLMsg     string
	TxtMsg      string
	Attachments []string
	Headers     map[string][]string
	Cc          []string
	Bcc         []string
	Recipients  []string
}

// InitializeMessage initializes the Message struct by creating sensible defaults for
// optional parameters that have zero values of nil.
func InitializeMessage() *Message {
	return &Message{
		Attachments: []string{},
		Headers:     make(map[string][]string),
		Cc:          []string{},
		Bcc:         []string{},
		Recipients:  []string{},
	}
}

// GetSubject ...
func (m *Message) GetSubject() string {
	return m.Subject
}

// GetHTMLMessage ...
func (m *Message) GetHTMLMessage() string {
	return m.HTMLMsg
}

// GetTextMessage ...
func (m *Message) GetTextMessage() string {
	return m.TxtMsg
}

// GetAttachments ...
func (m *Message) GetAttachments() []string {
	return m.Attachments
}

// GetHeaders ...
func (m *Message) GetHeaders() map[string][]string {
	return m.Headers
}

// GetCc ...
func (m *Message) GetCc() []string {
	return m.Cc
}

// GetBCc ...
func (m *Message) GetBCc() []string {
	return m.Bcc
}

// GetRecipients ...
func (m *Message) GetRecipients() []string {
	return m.Recipients
}

// ConsoleMailer is an implementation of the IMailer interface which writes mails
// to the console, it is suitable for debugging mails and for use in dev environments
type ConsoleMailer struct {
	mailer *Mailer
}

// NewConsoleMailer constructs an innstance of ConsoleMailer
func NewConsoleMailer() *ConsoleMailer {
	return &ConsoleMailer{&Mailer{}}
}

// Send builds emails from the provided messages and prints
// each email to the console
func (cm *ConsoleMailer) Send(msgs ...SendableMessage) error {
	messages := cm.mailer.buildMessages(msgs...)
	for _, msg := range messages {
		strBuilder := new(strings.Builder)
		msg.WriteTo(strBuilder)
		email := strBuilder.String()
		fmt.Println("================= BEGIN EMAIL =====================")
		fmt.Println("    ")
		fmt.Println(email)
		fmt.Println("    ")
		fmt.Println("================= END EMAIL =====================")
	}
	return nil
}
