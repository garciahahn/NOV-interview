package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	art "pkgs/artifact"
	"syscall"
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

// Statement to insert a new record in the raw data data table
const sqlStatementRaw = `
INSERT INTO raw (sensor, timestamp, value)
VALUES ($1, $2, $3)
RETURNING id`

// Statement to insert a new record in the averaged data data table
const sqlStatementAvg = `
INSERT INTO average (avg_ids, timestamp, value)
VALUES ($1, $2, $3)
RETURNING id`

func init() {
	// Running any external initialization routines
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Make signal channel to gracefully terminate program
	c := make(chan os.Signal, 32)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
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

	// deferring closing the connection when the program is finished
	defer db.Close()

	// Channel subscriber
	ch := make(chan *[]art.Sensor)
	_, err = ec.BindRecvChan(art.SensorChannel, ch)
	checkError(err)

	var id int = 0
	// Main working loop
	OuterLoop:
	for{
		select{
		// To gracefully break out of the main loop
		case <-c:
			fmt.Println("Program finished!")
			break OuterLoop
		default:
			// Wait for sensor data
			sensors := <- ch

			// Make arrays to store the table records and sensor values
			sensorValues := make([]float64, len(*sensors))
			tableIds := make([]int, len(*sensors))

			// Process stored data
			for i, s := range *sensors{
				// Writing in raw data table
				err = db.QueryRow(sqlStatementRaw,
								s.Name,
								s.Timestamp,
								s.Value).Scan(&id)
				checkError(err)
				tableIds[i] = id
				fmt.Printf("Added record %d from sensor %s with value %.2f\n",
							id,
							s.Name,
							s.Value)
				sensorValues[i] = s.Value
				}
		

			// Recording averaged values
			sensorAverage := avg(sensorValues)

			// Writing the averaged values in the average data table
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
}

// Small utility functions
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