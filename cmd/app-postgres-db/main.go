package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

//HealthCheck healthcheck endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Ok")
}

func main() {
	router := httprouter.New()
	router.GET("/", HealthCheck)
	go CreateData()
	log.Fatal(http.ListenAndServe(":5000", router))
}

//CreateData create data in the DB.
func CreateData() {
	// read configuration info from the environment variables
	dbUser := os.Getenv("PGUSER")
	dbPass := os.Getenv("PGPASSWORD")
	dbHost := os.Getenv("PGHOST")
	dbPort := os.Getenv("PGPORT")
	dbDatabase := os.Getenv("PGDATABASE")

	// open a connection to the table
	log.Printf("Connecting to postgres instance %s as user %s\n", dbHost, dbUser)
	connstring := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbDatabase)
	db, err := sql.Open("postgres", connstring)
	if err != nil {
		log.Println(`Could not connect to db`)
		panic(err)
	}
	defer db.Close()

	tableName := "app_postgres_db"

	// create the table if it doesn't already exist
	createstmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (Counter int);", tableName)
	log.Printf("Creating database table %s (if not exists)", tableName)
	_, err = db.Exec(createstmt)
	if err != nil {
		panic(err)
	}

	// find the current biggest value in the table
	var max int
	selectMax := fmt.Sprintf("SELECT MAX(Counter) FROM %s;", tableName)
	err = db.QueryRow(selectMax).Scan(&max)
	if err != nil {
		log.Printf("ignoring select max error: %+v\n", err)
	}

	log.Printf("MAX value in table is %d\n", max)

	// start inserting records into the table, starting with 1 bigger than the current max
	for index := 1; index <= 50; index++ {
		newRecord := max + index

		log.Printf("Inserting record %d in the database\n", newRecord)

		insertstmt := fmt.Sprintf("INSERT INTO %s values (%d);", tableName, newRecord)
		_, err = db.Exec(insertstmt)
		if err != nil {
			panic(err)
		}

		log.Println("Retrieving all records from the database")
		selectstmt := fmt.Sprintf("SELECT * FROM %s;", tableName)

		rows, err := db.Query(selectstmt)
		if err != nil {
			panic(err)
		}

		first := true
		var results strings.Builder
		var r int
		for rows.Next() {
			if !first {
				results.WriteString(", ")
			}

			if err := rows.Scan(&r); err != nil {
				log.Printf("Failed to scan next row, bailing out: %+v\n", err)
				break
			}

			results.WriteString(strconv.Itoa(r))
			first = false
		}

		log.Printf("All retrieved records: %s\n", results.String())
		time.Sleep(6 * time.Second)
	}

	log.Println("Completed inserting records into the database")
}
