package sendmail

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/pkg/errors"
)

// SendMail type
type SendMail struct {
	toAddress   map[string]string // Required *
	fromAddress string            // Required *
	subject     string            // Required *
	body        string            // Required *
	auth        smtp.Auth
	headers     map[string]string
	servername  string // Required *
	// attachments map[string][]byte
}

// Create a new SendMail
func New(servername, from, subject, body string) *SendMail {
	return &SendMail{
		fromAddress: from,
		subject:     subject,
		body:        body,
		servername:  servername,
	}
}

// Add an email address and name to, toAddress field
func (s *SendMail) AddToAddress(name string, email string) (bool, error) {
	curLen := len(s.toAddress)
	if curLen == 0 {
		s.toAddress = make(map[string]string)
	}

	if em, ok := s.toAddress[name]; ok && em != email { // Check if to address already exists
		return false, errors.New("To address already exists")
	}
	s.toAddress[name] = email
	newLen := len(s.toAddress)

	return newLen > curLen, nil
}

// Check if all the required fields have a value
func (s *SendMail) validate() error {
	var err error
	if l := len(s.toAddress); l == 0 {
		err = errors.New("Missing to address")
	} else if s.fromAddress == "" {
		err = errors.New("Missing from address")
	} else if s.subject == "" {
		err = errors.New("Missing subject")
	} else if s.body == "" {
		err = errors.New("Missing body")
	}

	return err
}

// Add a header to the mail header
func (s *SendMail) AddHeader(name, value string) {
	if len(s.headers) == 0 {
		s.headers = make(map[string]string)
	}

	s.headers[name] = value
}

// Add auth to send mail
func (s *SendMail) Auth(username, password, host string) {
	s.auth = smtp.PlainAuth("", username, password, host)
}

// Send email using the stmp.Dial
func (s *SendMail) Send() (bool, error) {
	var err error

	err = s.validate()
	if err != nil {
		return false, errors.Wrap(err, "Failed to send mail")
	}

	s.AddHeader("From", s.fromAddress)

	// Convert toAddress map to sendmail to address string
	var toAddressesHeader = ""
	var toAddresses []string
	for name, address := range s.toAddress {
		toAddressesHeader += name + " <" + address + ">, "
		toAddresses = append(toAddresses, address)
	}
	toAddressesHeader = strings.TrimRight(toAddressesHeader, ", ")
	s.AddHeader("To", toAddressesHeader)

	s.AddHeader("Subject", s.subject)

	// Set up message
	var message = ""
	for k, v := range s.headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + s.body

	// Connect to the smtp server
	c, err := smtp.Dial(s.servername)
	if err != nil {
		return false, errors.Wrap(err, "Failed to connect to smtp server")
	}

	// Add to auth
	err = c.Auth(s.auth)
	if err != nil {
		return false, errors.Wrap(err, "Filed to add auth")
	}

	// Add to and from address
	var fromEmail = s.fromAddress
	if strings.Contains(fromEmail, " <") {
		fromEmail = strings.Split(s.fromAddress, " <")[1]
		fromEmail = strings.Replace(fromEmail, ">", "", -1)
	}
	// log.Println(fromEmail)
	err = c.Mail(fromEmail)
	if err != nil {
		return false, errors.Wrap(err, "Failed to add from email address")
	}

	for _, address := range toAddresses {
		err = c.Rcpt(address)
		if err != nil {
			return false, errors.Wrap(err, "Failed to add recipient to email")
		}
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return false, errors.Wrap(err, "Failed get mail writter")
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return false, errors.Wrap(err, "Failed to write message")
	}

	err = w.Close()
	if err != nil {
		return false, errors.Wrao(err, "Failed to close writter")
	}

	err = c.Quit()
	if err != nil {
		return false, errors.Wrap(err, "Failed to close connection to mail server")
	}

	return true, err
}
