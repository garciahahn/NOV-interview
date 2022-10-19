package artifact

import (
	"math/rand"
	"time"
)

var SensorChannel string = "SENSOR_CHANNEL"

type Sensor struct {
	Name      string  `json:"name"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

func (s *Sensor) UpdateSensor(){
	s.Timestamp = time.Now().Unix()
	s.Value = rand.Float64() * 100
}

func NewSensor(name string) *Sensor{
	ret := new(Sensor)
	ret.Name = name
	ret.UpdateSensor()
	return ret
}