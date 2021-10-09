package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123DET123"
	dbname   = "gophercies_test"
)

func main() {
	f := float64(2)
	psqlInfo := fmt.Sprintf("host= %s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	handleError(err)
	handleError(resetDb(db, dbname))

	db.Close()
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	handleError(err)

	handleError(createPhoneNumbersTable(db))
	_, err = insertPhone(db, "1234567890")
	handleError(err)

	_, err = insertPhone(db, "1234567890")
	handleError(err)
	_, err = insertPhone(db, "123 456 7891")
	handleError(err)
	_, err = insertPhone(db, "(123) 456 7892")
	handleError(err)
	_, err = insertPhone(db, "(123) 456-7893")
	handleError(err)
	id, err := insertPhone(db, "123-456-7894")
	handleError(err)
	_, err = insertPhone(db, "123-456-7890")
	handleError(err)
	_, err = insertPhone(db, "1234567892")
	handleError(err)
	_, err = insertPhone(db, "(123)456-7892")
	handleError(err)

	phoneNumber, err := queryPhone(db, id)
	handleError(err)
	fmt.Println("Phone number is", phoneNumber)

	phones, err := getPhones(db)
	handleError(err)
	for _, p := range phones {
		fmt.Printf("Working on ....%+v\n", p)
		number := normalizePhoneNo(p.number)
		fmt.Println(number)
		if number == p.number {
			fmt.Println("No Change required")
		} else {
			existingPhone, _ := getPhoneByNumber(db, number)
			if existingPhone != -1 {
				fmt.Println("Deleting", p.number)
				handleError(deletePhone(db, p.id))

			} else {
				fmt.Println("Updating", p.number)
				handleError(updatePhone(db, p.id, number))
			}
		}
	}
	phonesUpdated, err := getPhones(db)
	fmt.Println("Updating")
	handleError(err)
	for _, p := range phonesUpdated {
		fmt.Printf("%+v\n", p)
	}
	defer db.Close()
}

type phoneNumber struct {
	id     int
	number string
}

func deletePhone(db *sql.DB, id int) error {
	statement := `DELETE FROM phone_numbers where id=$1`
	_, err := db.Exec(statement, id)
	return err
}

func getPhoneByNumber(db *sql.DB, number string) (int, error) {
	var id int
	err := db.QueryRow("select id from phone_numbers where value=$1", number).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func getPhones(db *sql.DB) ([]phoneNumber, error) {
	rows, err := db.Query("select id,value from phone_numbers order by value")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []phoneNumber
	for rows.Next() {
		var p phoneNumber
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func queryPhone(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("Select * from phone_numbers where id=$1", id).Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) returning id`

	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	} else {
		return id, nil
	}
}

func updatePhone(db *sql.DB, phoneId int, phoneNo string) error {
	statement := `UPDATE phone_numbers set value = $1 where id=$2`
	_, err := db.Exec(statement, phoneNo, phoneId)
	return err
}
func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers
	(id SERIAL, value VARCHAR)`
	_, err := db.Exec(statement)
	return err

}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
func resetDb(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	handleError(err)
	return createDb(db, name)
}
func createDb(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	handleError(err)
	return nil
}

func normalizePhoneNo(phone string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(phone, "")

}

// func normalizePhoneNo(phone string) string {
// 	var buf bytes.Buffer

// 	for _, ch := range phone {
// 		if ch >= '0' && ch <= '9' {
// 			buf.WriteRune(ch)
// 		}
// 	}

// 	return buf.String()
// }
