package rtd_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/aztechian/nextride-shortcut/internal/rtd"
	"github.com/stretchr/testify/assert"
)

func TestTripIsValid(t *testing.T) {
	arrival := time.Now().Add(3000 * time.Second).Unix()
	vehicleId := "1234"

	tests := []struct {
		name string
		trip rtd.Trip
		want bool
	}{
		{"Default object", rtd.Trip{}, false},
		{"Scheduled with no vehicle", rtd.Trip{Vehicle: &rtd.Vehicle{}, PredictedArrivalTime: &arrival, TripStopStatus: rtd.SCHEDULED}, true},
		{"Has vehicle, but no predicted arrival", rtd.Trip{Vehicle: &rtd.Vehicle{ID: vehicleId}}, true},
		{"CANCELLED trip", rtd.Trip{ScheduledArrivalTime: &arrival, TripStatus: rtd.CANCELLED}, false},
		{"Predicted arrival time, but no vehicle", rtd.Trip{PredictedArrivalTime: &arrival}, true},
		{"Valid object", rtd.Trip{Vehicle: &rtd.Vehicle{ID: vehicleId}, PredictedArrivalTime: &arrival, TripStopStatus: rtd.SCHEDULED}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.trip.IsValid())
		})
	}
}

func TestTripIsScheduled(t *testing.T) {
	t.Skip("Not implemented")
}

func TestTripIsShuttle(t *testing.T) {
	t.Skip("Not implemented")
}

func TestTripGetTime(t *testing.T) {
	arrivalTime := time.Now().Add(323*time.Second).Unix() * 1000 // RTD uses milliseconds
	expected := "5 minutes"
	trip := rtd.Trip{PredictedArrivalTime: &arrivalTime}
	assert.Equal(t, expected, trip.GetTime())
}

func TestGetUpcomingTrip(t *testing.T) {
	tests := []struct {
		name       string
		sourceFile string
		line       string
		wantNil    bool
	}{
		{"Matching Line, no valid trips", "../../test/25434.json", "E", true},
		{"Invalid line", "../../test/25434.json", "Z", true},
		{"Valid, duplicated line", "../../test/25434.json", "W", false},
		{"Valid Trips for R Line", "../../test/valid-vehicle-rline.json", "R", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := os.OpenFile(tt.sourceFile, os.O_RDONLY, 0)
			var station *rtd.Station
			_ = json.NewDecoder(data).Decode(&station)
			trip := station.GetUpcomingTrip(tt.line)
			if tt.wantNil {
				assert.Nil(t, trip)
			} else {
				assert.NotNil(t, trip)
			}
		})
	}
}
