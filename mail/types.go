package mail

import "gopkg.in/gomail.v2"

// SendableMessage interface is the interface accepted by mailer
// it contains details of the message
type SendableMessage interface {
	GetSubject() string
	GetHTMLMessage() string
	GetTextMessage() string
	GetAttachments() []string
	GetHeaders() map[string][]string
	GetCc() []string
	GetBCc() []string
	GetRecipients() []string
}

// IMailer the mail sender
type IMailer interface {
	Send(...SendableMessage) error
}

// IDialer is a low level implemntation detail of this package
// it describes an interface needed to send messages
type IDialer interface {
	DialAndSend(...*gomail.Message) error
}
