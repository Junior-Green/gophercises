package main

import (
	"fmt"
	"os"

	"github.com/Junior-Green/gophercises/normalize"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const databaseURL string = "postgres://postgres:password@localhost:5432/postgres"

func main() {
	db, err := sqlx.Open("pgx", databaseURL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	normalize.Init(db)

	numbers, err := normalize.GetAllPhoneNumbers(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	for _, num := range numbers {
		normalized := normalize.NormalizePhone(num)
		if err := normalize.DeleteNumber(db, normalized); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		if err := normalize.UpdateNumber(db, num, normalized); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
}
