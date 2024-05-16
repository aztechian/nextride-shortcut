package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aztechian/nextride-shortcut/internal/rtd"
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

func handler(w http.ResponseWriter, r *http.Request) {
	// Parse parameters from the URL query
	stopId := r.URL.Query().Get("stop")
	railLine := r.URL.Query().Get("line")

	// Check if both parameters are provided
	if stopId == "" || railLine == "" {
		http.Error(w, "Both stop and line parameters are required", http.StatusBadRequest)
		return
	}

	stopIdInt, err := strconv.ParseUint(stopId, 10, 64)
	if err != nil {
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
			http.Error(w, fmt.Sprintf("Failed to get upcoming trip: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "The Next %s Line arrives in %s", railLine, trip.GetTime())
	default:
		http.Error(w, fmt.Sprintf("%s is not a valid rail line", railLine), http.StatusBadRequest)
		return
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/next", http.StatusSeeOther)
}

func main() {
	// direct requests to root to /next
	http.HandleFunc("/", redirectHandler)

	// Define the route and handler function
	http.HandleFunc("/next", handler)

	// Start the HTTP server on port 8080
	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
