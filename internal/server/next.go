package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aztechian/nextride-shortcut/internal/rtd"
	"github.com/rs/zerolog/hlog"
)

const (
	A_LINE = "A"
	B_LINE = "B"
	D_LINE = "D"
	E_LINE = "E"
	G_LINE = "G"
	H_LINE = "H"
	L_LINE = "L"
	N_LINE = "N"
	R_LINE = "R"
	W_LINE = "W"
)

const (
	DEFAULT_LINE  = E_LINE
	UNION_STATION = "33727"
	ARAP_VILLAGE  = "34008"
	DEFAULT_STOP  = UNION_STATION
)

// NextHandler is an HTTP handler that returns the next train arrival time for a given stop and rail line
func NextHandler(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)
	// Parse parameters from the URL query
	stopId := r.URL.Query().Get("stop")
	railLine := r.URL.Query().Get("line")

	// Check if both parameters are provided
	if stopId == "" || railLine == "" {
		log.Debug().Str("stop", stopId).Str("line", railLine).Msg("Both stop and line parameters are required")
		http.Error(w, "Both stop and line parameters are required", http.StatusBadRequest)
		return
	}

	stopIdInt, err := strconv.ParseUint(stopId, 10, 64)
	if err != nil {
		log.Debug().Str("stop", stopId).Msg("Invalid stop ID, it must be a number")
		http.Error(w, "Invalid stop ID, it must be a number", http.StatusBadRequest)
		return
	}

	// Validate railLine value
	switch railLine {
	case A_LINE, B_LINE, D_LINE, E_LINE, G_LINE, H_LINE, L_LINE, N_LINE, R_LINE, W_LINE:

		// fmt.Fprintf(w, "Received stopId: %s and railLine: %s", stopId, railLine)
		api := rtd.NewRtd(railLine, stopIdInt)
		trip, err := api.Get()
		if err != nil {
			log.Debug().Err(err).Msg("Failed to get upcoming trip")
			http.Error(w, fmt.Sprintf("Failed to get upcoming trip: %v", err), http.StatusInternalServerError)
			return
		}
		if trip == nil {
			log.Debug().Str("stop", stopId).Str("line", railLine).Msg("No upcoming trips found")
			http.Error(w, fmt.Sprintf("No upcoming %s line trips found, check the RTD app", railLine), http.StatusNotFound)
			return
		}
		log.Debug().Str("stop", stopId).Str("line", railLine).Str("time", trip.GetTime()).Msg("Upcoming trip found and returned to user")
		fmt.Fprintf(w, "The Next %s Line arrives in %s", railLine, trip.GetTime())
	default:
		log.Debug().Str("line", railLine).Msg("Invalid rail line given")
		http.Error(w, fmt.Sprintf("%s is not a valid rail line", railLine), http.StatusBadRequest)
		return
	}
}
