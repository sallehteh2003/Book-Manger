package Validation

import (
	"Book_Manager/Database"
	"errors"
	"fmt"
)

var validDomain = []string{".com", ".ir", ".org"}

type Validation struct {
	db          *Database.GormDB
	ValidDomain []string
}

// Create a new instance of Validation

func CreateValidation(db *Database.GormDB) (*Validation, error) {
	if db == nil {
		return nil, errors.New("the database is essential for Validation")
	}
	return &Validation{
		db:          db,
		ValidDomain: validDomain,
	}, nil
}
func (v *Validation) ValidateData(email string, phone string) error {
	if err := v.validateEmail(email); err != nil {
		return err
	}
	if err := v.validatePhoneNumber(phone); err != nil {
		return err
	}
	return nil
}
func (v *Validation) validateEmail(email string) error {
	EmailChars := []rune(email)
	var Edesign = 0
	var Dot = false
	var temp int
	Domain := ""
	for i := 0; i < len(EmailChars); i++ {
		char := string(EmailChars[i])
		if char == "@" {
			Edesign++
			temp = i
		}
		if char == "." {
			if temp+1 == i {
				return errors.New("invalid email")
			}
			Dot = true

		}

	}
	fmt.Println("d")
	for i := len(EmailChars) - 1; i > -1; i-- {
		char := string(EmailChars[i])
		Domain = char + Domain
		if char == "." {
			break
		}
		if EmailChars[i] <= 58 && EmailChars[i] >= 48 {
			return errors.New("invalid email")
		}

	}
	if err := v.checkDomain(Domain); err != nil {
		return errors.New("invalid email")
	}
	if Edesign != 1 || !Dot {
		return errors.New("invalid email")
	}
	return nil
}
func (v *Validation) validatePhoneNumber(phone string) error {
	if len(phone) > 11 {
		return errors.New("invalid phone number")
	}
	for i := 0; i < len(phone); i++ {

		if phone[i] > 58 || phone[i] < 48 {
			return errors.New("invalid phone number")
		}

	}
	return nil
}

func (v *Validation) checkDomain(Domain string) error {
	for _, Dom := range v.ValidDomain {
		if Dom == Domain {
			return nil
		}
	}
	return errors.New("domain not found")
}
