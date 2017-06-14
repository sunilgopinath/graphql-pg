package zero

import (
	"database/sql"

	"github.com/lib/pq"
)

// Person models database user
type Person struct {
	ID        int
	FirstName string
	LastName  string
	Username  string
	Email     string
}

//StartDB returns temporary tables
func StartDB(db *sql.DB) error {
	query := `CREATE TEMP TABLE "userinfo" (
						  "uid" SERIAL PRIMARY KEY,
							"first_name" varchar(30),
							"last_name" varchar(30),
						  "username" varchar(30) UNIQUE NOT NULL,
							"email" varchar(30)
						)`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func LoadData(db *sql.DB) error {
	users := []Person{
		Person{
			FirstName: "Test", LastName: "User", Username: "demo1", Email: "test@fb.com",
		},
		Person{
			FirstName: "New", LastName: "User", Username: "demo2", Email: "new@fb.com",
		},
		Person{
			FirstName: "Latest", LastName: "User", Username: "demo3", Email: "latest@fb.com",
		},
	}
	txn, err := db.Begin()
	if err != nil {
		CheckErr(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("userinfo", "first_name", "last_name", "username", "email"))
	if err != nil {
		CheckErr(err)
	}

	for _, user := range users {
		_, err = stmt.Exec(user.FirstName, user.LastName, user.Username, user.Email)
		if err != nil {
			CheckErr(err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		CheckErr(err)
	}

	err = stmt.Close()
	if err != nil {
		CheckErr(err)
	}

	err = txn.Commit()
	if err != nil {
		CheckErr(err)
	}
	if err != nil {
		return err
	}
	return nil
}

//CheckErr returns a basic error if one exists
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
