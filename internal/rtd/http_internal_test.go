package rtd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		stationName string
		branchCount int
	}{
		{"valid-vehicle-rline", "../../test/valid-vehicle-rline.json", "Iliff Station", 2},
		{"shuttles", "../../test/colorado-shuttle-eline.json", "Colorado Station", 3},
		{"arapahoe", "../../test/cancelled-arapahoe.json", "Arapahoe at Village Center Station", 2},
	}

	for _, tt := range tests {
		// Read the test JSON file
		file, err := os.Open(tt.filename)
		if err != nil {
			t.Fatalf("%s: failed to open file: %v", tt.name, err)
		}
		defer file.Close()

		// Unmarshal the JSON data into a struct
		station, err := parseStation(file)
		if err != nil {
			t.Fatalf("%s: failed to unmarshal JSON: %v", tt.name, err)
		}

		assert.Equal(t, tt.stationName, station.Name)
		assert.Len(t, station.Branches, tt.branchCount)
	}
}
