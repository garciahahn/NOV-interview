package artifact

import (
	"math/rand"
	"time"
)

// Sensor channel name for NATS
const SensorChannel string = "SENSOR_CHANNEL"

// Sensor structure for encoding
type Sensor struct {
	Name      string  `json:"name"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// Update sensor with random data
func (s *Sensor) UpdateSensor(){
	s.Timestamp = time.Now().Unix()
	s.Value = rand.Float64() * 100
}

// Constuctor for a new sensor
func NewSensor(name string) *Sensor{
	ret := new(Sensor)
	ret.Name = name
	ret.UpdateSensor()
	return ret
}