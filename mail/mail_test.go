package mail

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dino16m/golearn-core/mail/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/gomail.v2"
)

type dummyMessage struct {
	subject     string
	htmlMsg     string
	txtMsg      string
	attachments []string
	headers     map[string][]string
	cc          []string
	bcc         []string
	recipients  []string
}

func initializeMsg() *dummyMessage {
	return &dummyMessage{
		attachments: []string{},
		headers:     make(map[string][]string),
		cc:          []string{},
		bcc:         []string{},
		recipients:  []string{},
	}
}

func (m *dummyMessage) GetSubject() string {
	return m.subject
}
func (m *dummyMessage) GetHTMLMessage() string {
	return m.htmlMsg
}
func (m *dummyMessage) GetTextMessage() string {
	return m.txtMsg
}
func (m *dummyMessage) GetAttachments() []string {
	return m.attachments
}
func (m *dummyMessage) GetHeaders() map[string][]string {
	return m.headers
}
func (m *dummyMessage) GetCc() []string {
	return m.cc
}
func (m *dummyMessage) GetBCc() []string {
	return m.bcc
}
func (m *dummyMessage) GetRecipients() []string {
	return m.recipients
}

type mailerTestSuite struct {
	suite.Suite
	dialer         mocks.IDialer
	senderName     string
	senderEmail    string
	mailer         *Mailer
	renderedEmails []string
}

func (s *mailerTestSuite) SetupTest() {
	dialer := new(mocks.IDialer)
	s.senderName = "dummy"
	s.senderEmail = "root@dummy.com"
	s.mailer = &Mailer{
		dialer:      dialer,
		senderName:  s.senderName,
		senderEmail: s.senderEmail,
	}
	s.renderedEmails = []string{}
	dialer.On("DialAndSend", mock.MatchedBy(
		func(msg *gomail.Message) bool {
			strBuilder := new(strings.Builder)
			msg.WriteTo(strBuilder)
			email := strBuilder.String()
			s.renderedEmails = append(s.renderedEmails, email)
			return true
		})).Return(nil)
}

func (s *mailerTestSuite) TestMailContainsCorrectSender() {
	msg := initializeMsg()
	msg.htmlMsg = "<h1>Hello I am me</h1>"
	msg.recipients = append(msg.recipients, "me@me.com")
	s.mailer.Send(msg)
	sender := fmt.Sprintf("\"%s\" <%s>", s.senderName, s.senderEmail)
	s.Contains(s.renderedEmails[0], sender)
}
func (s *mailerTestSuite) TestMailContainsCorrectRecipients() {
	msg := initializeMsg()
	msg.htmlMsg = "<h1>Hello I am me</h1>"
	recipient := "me@me.com"
	msg.recipients = append(msg.recipients, recipient)
	s.mailer.Send(msg)
	s.Contains(s.renderedEmails[0], recipient)
}
func (s *mailerTestSuite) TestMailContainsHTMLAndTextWhenSentTogether() {
	msg := initializeMsg()
	msg.htmlMsg = "<h1>Hello I am me</h1>"
	msg.txtMsg = "This is a text message"
	recipient := "me@me.com"
	msg.recipients = append(msg.recipients, recipient)
	s.mailer.Send(msg)
	email := s.renderedEmails[0]
	s.Contains(email, "text/html")
	s.Contains(email, msg.htmlMsg)
	s.Contains(email, "text/plain")
	s.Contains(email, msg.txtMsg)
}
func (s *mailerTestSuite) TestMailHasAttachmentsWhenSent() {
	msg := initializeMsg()
	msg.htmlMsg = "<h1>Hello I am me</h1>"
	msg.recipients = append(msg.recipients, "me@me.com")
	tmpDir, err := ioutil.TempDir("", "attachments")
	if err != nil {
		panic(err)
	}
	attachmentFile := "attachment.txt"
	attachmentFilePath := filepath.Join(tmpDir, attachmentFile)
	ioutil.WriteFile(attachmentFilePath, []byte("This is an attachment"), 0666)
	msg.attachments = append(msg.attachments, attachmentFilePath)
	s.mailer.Send(msg)
	email := s.renderedEmails[0]
	s.Contains(email, fmt.Sprintf("filename=\"%s\"", attachmentFile))
	os.Remove(tmpDir)
}
func (s *mailerTestSuite) TestHeadersAllSet() {
	msg := initializeMsg()
	msg.htmlMsg = "<h1>Hello I am me</h1>"
	msg.recipients = append(msg.recipients, "me@me.com")
	subject := "test"
	msg.headers["Subject"] = []string{subject}
	s.mailer.Send(msg)
	email := s.renderedEmails[0]
	s.Contains(email, fmt.Sprintf("Subject: %s", subject))
}

func (s *mailerTestSuite) TestSubjectSet() {
	msg := initializeMsg()
	msg.htmlMsg = "<h1>Hello I am me</h1>"
	msg.recipients = append(msg.recipients, "me@me.com")
	subject := "test"
	msg.subject = subject
	s.mailer.Send(msg)
	email := s.renderedEmails[0]
	s.Contains(email, fmt.Sprintf("Subject: %s", subject))
}
func (s *mailerTestSuite) TestBccSet() {
	msg := initializeMsg()
	msg.htmlMsg = "<h1>Hello I am me</h1>"
	msg.recipients = append(msg.recipients, "me@me.com")
	bcc := "bcc@hotmail.com"
	msg.bcc = append(msg.bcc, bcc)
	s.mailer.Send(msg)
	email := s.renderedEmails[0]
	s.Contains(email, fmt.Sprintf("BCc: %s", bcc))
}

func (s *mailerTestSuite) TestCcSet() {
	msg := initializeMsg()
	msg.htmlMsg = "<h1>Hello I am me</h1>"
	msg.recipients = append(msg.recipients, "me@me.com")
	cc := "bcc@hotmail.com"
	msg.cc = append(msg.cc, cc)
	s.mailer.Send(msg)
	email := s.renderedEmails[0]
	s.Contains(email, fmt.Sprintf("Cc: %s", cc))
}
func (s *mailerTestSuite) TestMailContainsOnlyTextWhenOnlyTextSet() {
	msg := initializeMsg()
	msg.txtMsg = "This is a text message"
	recipient := "me@me.com"
	msg.recipients = append(msg.recipients, recipient)
	s.mailer.Send(msg)
	email := s.renderedEmails[0]
	s.Contains(email, "text/plain")
	s.Contains(email, msg.txtMsg)
	s.NotContains(email, "text/html")
}
func (s *mailerTestSuite) TestMailContainsOnlyHTMLWhenOnlyHTMLSet() {
	msg := initializeMsg()
	msg.htmlMsg = "This is a html message"
	recipient := "me@me.com"
	msg.recipients = append(msg.recipients, recipient)
	s.mailer.Send(msg)
	email := s.renderedEmails[0]
	s.Contains(email, "text/html")
	s.Contains(email, msg.htmlMsg)
	s.NotContains(email, "text/plain")
}
func TestMailer(t *testing.T) {
	suite.Run(t, new(mailerTestSuite))
}
