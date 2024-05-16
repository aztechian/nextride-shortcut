package rtd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Rtd struct {
	Client *http.Client
	Line   string
	Stop   uint64
}

const (
	URL_TEMPLATE = "https://nodejs-prod.rtd-denver.com/api/v2/nextride/stops/%d"
	API_KEY      = "e7b926a1-cddb-46e7-bb27-6d134e5b5feb"
	BROWSER      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15"
)

// NewRtd creates a new Rtd struct with a sensible default http.Client (10 second timeout)
func NewRtd(line string, stop uint64) *Rtd {
	return &Rtd{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Line: line,
		Stop: stop,
	}
}

func NewRequest(line string, stop uint64) (*http.Request, error) {
	url := fmt.Sprintf(URL_TEMPLATE, stop)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	setHeaders(req)
	return req, nil
}

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-key", API_KEY)
	req.Header.Set("User-Agent", BROWSER)
	req.Header.Set("Referer", "https://app.rtd-denver.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
}

func parseStation(body io.ReadCloser) (*Station, error) {
	var station *Station
	err := json.NewDecoder(body).Decode(&station)
	if err != nil {
		return nil, err
	}
	return station, nil
}

func (r Rtd) DoRequest(req *http.Request) (*http.Response, error) {
	setHeaders(req)
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Rtd) Get() (*UpcomingTrip, error) {
	req, err := NewRequest(r.Line, r.Stop)
	if err != nil {
		return nil, err
	}

	resp, err := r.DoRequest(req)
	if err != nil {
		return nil, err
	}

	body := resp.Body
	station, err := parseStation(body)
	if err != nil {
		return nil, err
	}

	trip := station.GetUpcomingTrip(r.Line)
	return trip, nil
}
