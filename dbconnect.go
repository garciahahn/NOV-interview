package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

const(
	host = "localhost"
	port = 5432
	user = "NOV"
	password = "12346"
	dbname = "sensor_data"
)

const sqlStatement = `
INSERT INTO raw (sensor, timestamp, value)
VALUES ($1, $2, $3)
RETURNING id, value`



func main() {
	rand.Seed(time.Now().UnixNano())
	psqInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
                           "dbname=%s sslmode=disable",
						   host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqInfo)
	handleError(err)
	defer db.Close()

	err = db.Ping()
	handleError(err)
	id := 0
	var value float64
	err = db.QueryRow(sqlStatement,
					  "Sensor 1",
					  time.Now().Unix(),
					  rand.Float64()*100).
					  Scan(&id, &value)

	handleError(err)
	fmt.Println("New record ID is:", id, "and has value", value)
	fmt.Println("Successfully connected!")
}

func handleError(err error){
	if err != nil{
		panic(err)
	}
}