package rtd

import (
	"fmt"
	"time"
)

const (
	SCHEDULED = "SCHEDULED"
	SHUTTLE   = "BUS SHUTTLE"
	CANCELLED = "CANCELLED"
	MS_TO_SEC = 1000
)

type Station struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Lat               float64   `json:"lat"`
	Lng               float64   `json:"lng"`
	ChildStops        []string  `json:"childStops"`
	ParentStationName string    `json:"parentStationName"`
	ParentStationID   string    `json:"parentStationId"`
	Branches          []*Branch `json:"branches"`
}

type Branch struct {
	ID             string  `json:"id"`
	RouteColor     string  `json:"routeColor"`
	RouteTextColor string  `json:"routeTextColor"`
	RouteLongName  string  `json:"routeLongName"`
	RouteType      *int    `json:"routeType"`
	Mode           string  `json:"mode"`
	Headsign       string  `json:"headsign"`
	DirectionID    *int    `json:"directionId"`
	DirectionName  string  `json:"directionName"`
	StopName       string  `json:"stopName"`
	StopID         string  `json:"stopId"`
	AgencyID       string  `json:"agencyId"`
	DropoffOnly    bool    `json:"dropoffOnly"`
	UpcomingTrips  []*Trip `json:"upcomingTrips"`
}

type Trip struct {
	PredictedArrivalTime   *int64   `json:"predictedArrivalTime"`
	PredictedDepartureTime *int64   `json:"predictedDepartureTime"`
	ScheduledArrivalTime   *int64   `json:"scheduledArrivalTime"`
	ScheduledDepartureTime *int64   `json:"scheduledDepartureTime"`
	StopDropOffType        *int     `json:"stopDropOffType"`
	StopPickupType         *int     `json:"stopPickupType"`
	TripID                 string   `json:"tripId"`
	TripStatus             string   `json:"tripStatus"`
	TripStopStatus         string   `json:"tripStopStatus"`
	Vehicle                *Vehicle `json:"vehicle,omitempty"`
}

type Vehicle struct {
	Bearing        *float32 `json:"bearing"`
	ID             string   `json:"id"`
	Label          string   `json:"label"`
	Lat            *float32 `json:"lat"`
	Lng            *float32 `json:"lng"`
	Timestamp      *int64   `json:"timestamp"`
	RouteTextColor string   `json:"routeTextColor"`
	RouteColor     string   `json:"routeColor"`
	Mode           string   `json:"mode"`
}

func (t *Trip) IsValid() bool {
	return t.IsScheduled() || t.IsPredicted() || t.HasVehicle()
}

func (t *Trip) IsScheduled() bool {
	return t.TripStopStatus == SCHEDULED
}

func (t *Trip) IsShuttleBus() bool {
	return t.TripStopStatus == SHUTTLE
}

func (t *Trip) IsCancelled() bool {
	return t.TripStatus == CANCELLED
}

func (t *Trip) IsPredicted() bool {
	return t.PredictedArrivalTime != nil
}

func (t *Trip) HasVehicle() bool {
	return t.Vehicle != nil && t.Vehicle.ID != ""
}

func (t *Trip) GetTime() string {
	if t.PredictedArrivalTime != nil {
		return relativeTime(*t.PredictedArrivalTime)
	} else if t.ScheduledArrivalTime != nil {
		return relativeTime(*t.ScheduledArrivalTime)
	}
	return "unknown"
}

func relativeTime(t int64) string {
	targetTime := time.Unix(t/MS_TO_SEC, 0)
	duration := time.Until(targetTime).Minutes()

	return fmt.Sprintf("%d minutes", int(duration))
}

func (s *Station) GetUpcomingTrip(line string) *Trip {
	for _, branch := range s.Branches {
		if branch.ID == line {
			for _, trip := range branch.UpcomingTrips {
				if trip.IsValid() {
					return trip
				}
			}
		}
	}
	return nil
}
