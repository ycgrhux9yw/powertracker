package client

import (
	"time"
)

// PowerReading represents a single power consumption reading from the meter.
type PowerReading struct {
	// Timestamp is when the reading was taken.
	Timestamp time.Time `json:"timestamp"`

	// WattsNow is the current power consumption in watts.
	WattsNow float64 `json:"watts_now"`

	// WattHoursToday is the cumulative energy consumed today in watt-hours.
	WattHoursToday float64 `json:"watt_hours_today"`

	// WattHoursTotal is the total cumulative energy consumed in watt-hours.
	WattHoursTotal float64 `json:"watt_hours_total"`

	// Voltage is the current voltage reading, if available.
	Voltage float64 `json:"voltage,omitempty"`

	// Source identifies which meter or channel produced this reading.
	Source string `json:"source"`
}

// MetricsCollector aggregates power readings and provides summary statistics.
type MetricsCollector struct {
	readings []PowerReading
	maxSize  int
}

// NewMetricsCollector creates a MetricsCollector that retains up to maxSize readings.
func NewMetricsCollector(maxSize int) *MetricsCollector {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &MetricsCollector{
		readings: make([]PowerReading, 0, maxSize),
		maxSize:  maxSize,
	}
}

// Add appends a new reading to the collector, evicting the oldest if at capacity.
func (m *MetricsCollector) Add(r PowerReading) {
	if len(m.readings) >= m.maxSize {
		// Shift slice to drop oldest entry.
		m.readings = m.readings[1:]
	}
	if r.Timestamp.IsZero() {
		r.Timestamp = time.Now()
	}
	m.readings = append(m.readings, r)
}

// Latest returns the most recent reading, and false if no readings exist.
func (m *MetricsCollector) Latest() (PowerReading, bool) {
	if len(m.readings) == 0 {
		return PowerReading{}, false
	}
	return m.readings[len(m.readings)-1], true
}

// AverageWatts returns the mean power consumption across all retained readings.
// Returns 0 if there are no readings.
func (m *MetricsCollector) AverageWatts() float64 {
	if len(m.readings) == 0 {
		return 0
	}
	var sum float64
	for _, r := range m.readings {
		sum += r.WattsNow
	}
	return sum / float64(len(m.readings))
}

// Count returns the number of readings currently held.
func (m *MetricsCollector) Count() int {
	return len(m.readings)
}

// Reset clears all stored readings.
func (m *MetricsCollector) Reset() {
	m.readings = m.readings[:0]
}
