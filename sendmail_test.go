package sendmail

import "testing"

func TestAddToAddress(t *testing.T) {
	mail := New("", "", "", "")
	ok, err := mail.AddToAddress("Test Person", "test_email@test.mail")
	if err != nil {
		t.Error(err.Error())
	}

	if ok == false {
		t.Error("SendMail.AddToAddress failed")
	}

	// Add another email address
	ok, err = mail.AddToAddress("Test Again", "tes@mail.testing")
	if err != nil {
		t.Error(err.Error())
	}

	if ok == false {
		t.Error("SendMail.AddToAddress failed")
	}

	l := len(mail.toAddress)
	if l != 2 {
		t.Error("Second email address wasn't added to SendMail")
	}

	// Add an email address that already exists
	ok, err = mail.AddToAddress("Test Again", "new@address.email")
	if err == nil || err.Error() != "To address already exists" {
		t.Error("Failed and check for existsing to name")
	}

	if ok == true {
		t.Error("Shouldn't have added new email address to already existing name key")
	}
}

func TestAuth(t *testing.T) {
	mail := New("", "", "", "")
	mail.Auth("username", "password", "hostname / ip")
}

func TestAddHeader(t *testing.T) {
	mail := New("", "", "", "")
	mail.AddHeader("Content-Type", "text/html")

	if mail.headers["Content-Type"] != "text/html" {
		t.Error("Should have added content type to header")
	}
}

func TestValidate(t *testing.T) {
	// TODO: Complete test with auth and servername
	mail := New("", "", "", "")
	err := mail.validate()
	if err == nil {
		t.Error("Failed to validate to address")
	}

	// Add to address and check for valid from address
	mail.AddToAddress("test", "test")
	err = mail.validate()
	if err == nil {
		t.Error("Failed to validate to from address")
	}

	// Add from address and check for valid subject
	mail = New("", "test@te.st", "", "")
	mail.AddToAddress("test", "test")
	err = mail.validate()
	if err == nil {
		t.Error("Faild to validate subject")
	}

	// Add subject and check for valid body
	mail = New("", "test@te.st", "Test Subject", "")
	mail.AddToAddress("test", "test")
	err = mail.validate()
	if err == nil {
		t.Error("Failed to validate body")
	}
}

func TestSend(t *testing.T) {
	mail := New(
		"(hostname / ip) :port",
		"fromAddr",
		"subject",
		"Body content",
	)

	mail.AddToAddress("Test", "test@test.com")
	mail.AddToAddress("Test2", "test2@te.st")

	mail.Auth("test@te.st", "passwd", "hostname / ip")

	_, err := mail.Send()
	if err != nil {
		t.Error(err)
	}
}
