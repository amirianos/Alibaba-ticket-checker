package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	Result []struct {
		Price           float64 `json:"price"`
		IsAvailable     bool    `json:"isAvailable"`
		IsSoldOut       bool    `json:"isSoldOut"`
		IsCurrentDate   bool    `json:"isCurrentDate"`
		IsMostExpensive bool    `json:"isMostExpensive"`
		IsBestPrice     bool    `json:"isBestPrice"`
		Date            string  `json:"date"`
	} `json:"result"`
	TargetURL           interface{} `json:"targetUrl"`
	Success             bool        `json:"success"`
	Error               interface{} `json:"error"`
	UnauthorizedRequest bool        `json:"unauthorizedRequest"`
	Wrapped             bool        `json:"__wrapped"`
	TraceID             interface{} `json:"__traceId"`
}

type Config struct {
	Date        string `json:"date"`
	Destination string `json:"destination"`
	Origin      string `json:"origin"`
}

func SendRequest(url string) (*Request, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var request Request
	err = json.NewDecoder(resp.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}
func urlMaker(config *Config) string {
	url := "https://ws.alibaba.ir/api/v2/train/price-calender?Origin=" + config.Origin + "&Destination=" + config.Destination + "&Date=" + config.Date + "&CountDay=15&TicketType=Family"
	return url
}
func Readconfigfile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
func main() {
	configfile, _ := Readconfigfile("config.json")
	url := urlMaker(configfile)
	sentrequest, _ := SendRequest(url)
	date := configfile.Date
	sample_date := "2006-01-02"
	timedate, _ := time.Parse(sample_date, date)
	utc_format := timedate.UTC().Format("2006-01-02T00:00:00")
	ticket_availibility := false
	for _, value := range sentrequest.Result {
		if value.Date == utc_format {
			ticket_availibility = true
			if !value.IsAvailable {
				fmt.Println("Currently tickets alls are reserved!")
			} else {
				fmt.Println("Ticket found")
				fmt.Println("Ticket price is:", value.Price)
			}
		}
	}
	if !ticket_availibility {
		fmt.Println("Ticket does not found!")
	}
}
