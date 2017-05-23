package sendmail

import (
	"log"
	"testing"
)

func TestAddToAddress(t *testing.T) {
	mail := SendMail{}
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

	l := len(mail.ToAddress)
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

func TestValidate(t *testing.T) {
	mail := SendMail{}
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
	mail.FromAddress = "test@te.st"
	err = mail.validate()
	if err == nil {
		t.Error("Faild to validate subject")
	}

	// Add subject and check for valid body
	mail.Subject = "Test Subject"
	err = mail.validate()
	if err == nil {
		t.Error("Failed to validate body")
	}
}

func TestSend(t *testing.T) {
	mail := SendMail{
		FromAddress: "test@te.st",
		Subject:     "Test Subject",
		Body:        "Test body",
	}

	mail.AddToAddress("Test", "test@te.st")
	mail.AddToAddress("Test2", "test2@te.st")
	_, err := mail.Send()
	if err != nil {
		log.Println(err)
		t.Error(err)
	}
}
