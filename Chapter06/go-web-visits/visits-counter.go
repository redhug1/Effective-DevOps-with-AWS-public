package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/visits", viewHandler)
	fmt.Printf("doing ListenAndServe ...\n")
	err := http.ListenAndServe(":80", nil)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		log.Fatal()
	}
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	visitsCounter := updateDbCounter()
	templateText := "<h1>By Red:</h1>\n\nVisits Counter: {{.}}\n"
	tmpl, err := template.New("count").Parse(templateText)
	check(err)
	err = tmpl.Execute(writer, visitsCounter)
	check(err)
}

func updateDbCounter() int64 {
	db := dbConn()
	defer db.Close()
	_, err := db.Exec("UPDATE visits SET count = count+1 WHERE id = 1")
	check(err)

	var ID int64
	var count int64
	var version int64

	row := db.QueryRow("SELECT * FROM visits WHERE id = 1")
	if err = row.Scan(&ID, &count, &version); err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("no row\n")
			return -1
		}
		fmt.Printf("other error %v\n", err)
		return -2
	}

	return count
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "monty"
	dbPass := "some_pass"
	dbName := "demodb"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
