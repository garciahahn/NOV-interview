package main

import (
	art "pkgs/artifact"
	"time"

	"github.com/nats-io/nats.go"
)

var sleepTime = time.Second * 1

func main() {
	// Connect to server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	// Set up the json encoder for communication
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}
	defer ec.Close()

	sensors := initSensors()

	for{
		for i := 0; i<len(sensors); i++{
			sensors[i].UpdateSensor()
		}
		if err := ec.Publish(art.SensorChannel, sensors);err != nil{
			panic(err)
		}
		time.Sleep(sleepTime)
	}

}

// Initializate the three sensors, this could be further developed into a 'n'
// amount of sensors function
func initSensors() []*art.Sensor{
	return []*art.Sensor{art.NewSensor("Sensor 1"),
	art.NewSensor("Sensor 2"),
	art.NewSensor("Sensor 3")}
}