package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
  A_LINE = "a"
  B_LINE = "b"
  D_LINE = "d"
  E_LINE = "e"
  G_LINE = "g"
  H_LINE = "h"
  L_LINE = "l"
  N_LINE = "n"
  R_LINE = "r"
  W_LINE = "w"
)

const (
  DEFAULT_LINE = E_LINE
  UNION_STATION = "33727"
  ARAP_VILLAGE = "34008"
  DEFAULT_STOP = UNION_STATION
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Parse parameters from the URL query
	stopId := r.URL.Query().Get("stopId")
	railLine := r.URL.Query().Get("railLine")

	// Check if both parameters are provided
	if stopId == "" || railLine == "" {
		http.Error(w, "Both stopId and railLine parameters are required", http.StatusBadRequest)
		return
	}

	// Process the parameters
	// You can perform any operations here based on the received parameters

	// Respond to the client
	fmt.Fprintf(w, "Received stopId: %s and railLine: %s", stopId, railLine)
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
