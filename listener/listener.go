package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	art "pkgs/artifact"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

// Database connection information
const(
	host = "localhost"
	port = 5432
	user = "NOV"
	password = "12346"
	dbname = "sensor_data"
)

// Statement to insert a new record
const sqlStatementRaw = `
INSERT INTO raw (sensor, timestamp, value)
VALUES ($1, $2, $3)
RETURNING id`

const sqlStatementAvg = `
INSERT INTO average (avg_ids, timestamp, value)
VALUES ($1, $2, $3)
RETURNING id`

func init() {
	// Running any external initialization routines
	rand.Seed(time.Now().UnixNano())
}

func main() {

	// Connecting to the NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	checkError(err)

	// Encoding the connection with Json
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	checkError(err)

	// Connecting to the database
	psqInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
                           "dbname=%s sslmode=disable",
						   host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqInfo)
	checkError(err)
	// Pushing closing the connection when the program is finished
	defer db.Close()

	// Channel subscriber
	ch := make(chan *[]art.Sensor)
	_, err = ec.BindRecvChan(art.SensorChannel, ch)
	checkError(err)

	var id int = 0
	for{
		sensors := <- ch
		sensorValues := make([]float64, len(*sensors))
		tableIds := make([]int, len(*sensors))
		for i, s := range *sensors{
			err = db.QueryRow(sqlStatementRaw,
							  s.Name,
							  s.Timestamp,
							  s.Value).Scan(&id)
			checkError(err)
			tableIds[i] = id
			fmt.Printf("Added record %d from sensor %s\n", id, s.Name)
			sensorValues[i] = s.Value
		}

		// Recording averaged values
		sensorAverage := avg(sensorValues)
		err = db.QueryRow(sqlStatementAvg,
						pq.Array(tableIds),
						time.Now().Unix(),
						sensorAverage).Scan(&id)

		checkError(err)
		fmt.Printf("Added average record %d with value %.2f\n",
					id,
					sensorAverage)
	}
}

func sum(arr []float64) float64{
	var ret float64 = 0.0
	for _, v := range arr{
		ret += v
	}

	return ret
}

func avg(arr []float64) float64{
	return sum(arr) / float64(len(arr))
}

func checkError(err error){
	if err != nil{
		panic(err)
	}
}