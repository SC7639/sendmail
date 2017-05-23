package sendmail

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"
)

type SendMail struct {
	ToAddress   map[string]string // Required *
	FromAddress string            // Required *
	Subject     string            // Required *
	Body        string            // Required *
	// attachments map[string][]byte
}

// Add an email address and name to, toAddress field
func (s *SendMail) AddToAddress(name string, email string) (bool, error) {
	curLen := len(s.ToAddress)
	if curLen == 0 {
		s.ToAddress = make(map[string]string)
	}

	if em, ok := s.ToAddress[name]; ok && em != email { // Check if to address already exists
		return false, errors.New("To address already exists")
	}
	s.ToAddress[name] = email
	newLen := len(s.ToAddress)

	return newLen > curLen, nil
}

// Check if all the required fields have a value
func (s *SendMail) validate() error {
	var err error
	if l := len(s.ToAddress); l == 0 {
		err = errors.New("Missing to address")
	} else if s.FromAddress == "" {
		err = errors.New("Missing from address")
	} else if s.Subject == "" {
		err = errors.New("Missing subject")
	} else if s.Body == "" {
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
	for name, address := range s.ToAddress {
		toAddresses += name + " <" + address + ">, "
	}
	toAddresses = strings.TrimRight(toAddresses, ", ")

	// Create sendmail message
	var msg = "To: " + toAddresses
	msg += "Subject: " + s.Subject + "\n"
	msg += s.Body + "\n"

	sendmail := exec.Command("sendmail", "-f", s.FromAddress, toAddresses)
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
