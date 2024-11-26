package normalize

import (
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PhoneNumber struct {
	Number string `db:"phone_number"`
}

const phoneSchema = `
CREATE TABLE Phone (
    phone_number VARCHAR(128)
);`

func Init(db *sqlx.DB) {
	resetDB(db)
	populateDB(db)
}

func populateDB(db *sqlx.DB) {
	var phoneNumbers = [...]PhoneNumber{
		{"1234567890"},
		{"123 456 7891"},
		{"(123) 456 7892"},
		{"(123) 456-7893"},
		{"123-456-7894"},
		{"123-456-7890"},
		{"1234567892"},
		{"(123)456-7892"},
	}

	if _, err := db.NamedExec("INSERT INTO Phone (phone_number) VALUES (:phone_number)", phoneNumbers); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func NormalizePhone(phone PhoneNumber) PhoneNumber {
	return PhoneNumber{strings.Map(func(char rune) rune {
		if char >= 48 && char <= 57 {
			return char
		}
		return -1
	}, phone.Number)}
}

func resetDB(db *sqlx.DB) {
	if _, err := db.Query("SELECT * FROM Phone limit 1"); err == nil {
		db.MustExec(`DROP TABLE Phone`)
	}
	db.MustExec(phoneSchema)
}

func GetAllPhoneNumbers(db *sqlx.DB) ([]PhoneNumber, error) {
	var numbers []PhoneNumber
	if err := db.Select(&numbers, "SELECT * FROM Phone"); err != nil {
		return nil, err
	}
	return numbers, nil
}

func AddNumber(db *sqlx.DB, phoneNumber PhoneNumber) error {
	if _, err := db.NamedExec("INSERT INTO Phone (phone_number) VALUES (:phone_number)", phoneNumber); err != nil {
		return err
	}
	return nil
}

func DeleteNumber(db *sqlx.DB, phoneNumber PhoneNumber) error {
	if _, err := db.NamedExec("DELETE FROM Phone WHERE phone_number=(:phone_number)", phoneNumber); err != nil {
		return err
	}
	return nil
}

func UpdateNumber(db *sqlx.DB, phoneNumber PhoneNumber, to PhoneNumber) error {
	if _, err := db.Exec("UPDATE Phone SET phone_number = $1 WHERE phone_number = $2", to.Number, phoneNumber.Number); err != nil {
		return err
	}
	return nil
}
