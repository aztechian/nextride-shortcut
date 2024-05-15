import (
  "fmt"
  "net/http"
)

const(
	URL_TEMPLATE = "https://www.rtd-denver.com/api/nextride/stops/%d"
)

func NewRequest(line, stop) (http.Request, error) {
	url := fmt.Sprintf(URL_TEMPLATE, stop)
	
	req, err := http.NewRequest(http.MethodGet, url, 
