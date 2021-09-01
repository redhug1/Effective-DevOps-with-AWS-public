package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Print out counter to check that db is running
	count, result := updateDbCounter()
	fmt.Printf("visits counter is: %d, %s\n", count, result)

	http.HandleFunc("/visits", viewHandler)
	fmt.Printf("doing ListenAndServe ...\n")
	// port 80 is redirected to 8080 with iptables rules in terraform setup code
	err := http.ListenAndServe(":8080", nil)
	check(err)
}

func check(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		log.Fatal()
	}
}

type Info struct {
	Counter int64
	Status  string
	Ip      string
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
	visitsCounter, result := updateDbCounter()
	info := Info{visitsCounter, result, GetLocalIP()}

	templateText := "<h1>By Red:</h1>\n<p>Visits Counter: {{.Counter}}</p>\n<p>Status: {{.Status}}</p>\n<p>Private IP: {{.Ip}}</p>\n"
	t, err := template.New("count").Parse(templateText)
	check(err)
	err = t.Execute(writer, info)
	check(err)
}

func updateDbCounter() (int64, string) {
	db, err := dbConn()
	if err != nil {
		return -3, fmt.Sprintf("%v", err)
	}
	defer db.Close()
	_, err = db.Exec("UPDATE visits SET count = count+1 WHERE id = 1")
	if err != nil {
		return -4, fmt.Sprintf("%v", err)
	}

	var ID int64
	var count int64
	var version int64

	row := db.QueryRow("SELECT * FROM visits WHERE id = 1")
	if err = row.Scan(&ID, &count, &version); err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("no row\n")
			return -1, "no row"
		}
		fmt.Printf("other error %v\n", err)
		return -2, fmt.Sprintf("other error: %v", err)
	}

	fmt.Printf("count: %d\n", count)

	return count, "ok"
}

func dbConn() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "monty"
	dbPass := "some_pass"
	dbName := "demodb"
	// dbUrl := ""
	dbUrl := "tcp(localhost:3306)"
	// dbUrl := "tcp(<Endpoint field>:3306)"
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@"+dbUrl+"/"+dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
