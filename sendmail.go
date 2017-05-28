package sendmail

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
)

// SendMail type
type SendMail struct {
	toAddress   map[string]string // Required *
	fromAddress string            // Required *
	subject     string            // Required *
	body        string            // Required *
	// attachments map[string][]byte
}

// Create a new SendMail
func New(from, subject, body string) *SendMail {
	return &SendMail{
		fromAddress: from,
		subject:     subject,
		body:        body,
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

// Send email using the send mail command
func (s *SendMail) Send() (bool, error) {
	var success = false
	var err error

	err = s.validate()
	if err != nil {
		return success, err
	}

	// Convert toAddress map to sendmail to address string
	var toAddresses = ""
	for name, address := range s.toAddress {
		toAddresses += name + " <" + address + ">, "
	}
	toAddresses = strings.TrimRight(toAddresses, ", ")

	// Create sendmail message
	var msg = "To: " + toAddresses
	msg += "Subject: " + s.subject + "\n"
	msg += s.body + "\n"

	sendmail := exec.Command("sendmail", "-f", s.fromAddress, toAddresses)
	stdin, err := sendmail.StdinPipe()
	if err != nil {
		return success, err
	}

	stdout, err := sendmail.StdoutPipe() // Combine stdout and stderr
	if err != nil {
		return success, err
	}

	err = sendmail.Start() // Start sendmail command
	if err != nil {
		return success, err
	}

	_, err = stdin.Write([]byte(msg)) // Write message to sendmail
	if err != nil {
		stdin.Close()
		return success, err
	}
	stdin.Close()

	sentBytes, _ := ioutil.ReadAll(stdout)
	sendmail.Wait()
	if sentBytes == nil || string(sentBytes) == "" {
		success = true
	} else {
		err = errors.New(string(sentBytes))
	}

	return success, err
}
