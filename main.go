package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	API_PUSH = "/push"
)

func main() {
	filePath := flag.String("file", "", "path to a valid calender definition file in json format")
	token := flag.String("token", "", "users authentication token")
	host := flag.String("host", "http://localhost:5000", "operating server")
	flag.Parse()

	if *token == "" {
		log.Fatalln("token must be specified")
	}

	if *filePath != "" {
		calendars := loadFile(*filePath)
		push(calendars, *host)
		return
	}

	run(*token, *host)
}

func run(userToken string, host string) {
	var calendars []*Calendar
	for {
		date := readLine("Enter date [\"DD.MM.YYYY\"] (empty for today)")
		var dateAsUnix int64

		if date == "" {
			dateAsUnix = int64(time.Now().Unix()/60/60/24) * (60 * 60 * 24)
		} else {
			t, err := time.Parse("02.01.2006", date)
			if err != nil {
				fmt.Printf("error parsing date, invalid format: %s\n", err.Error())
				continue
			}
			dateAsUnix = t.Unix()
		}

		events := createEvents()
		calendars = append(calendars, &Calendar{
			UserToken: userToken,
			Date:      dateAsUnix,
			Events:    events,
		})

		selection := readLine("Push [1], Export json [2], Both [3], Add another Calendar [any]: ")
		if selection == "1" || selection == "3" {
			push(calendars, host)
			break
		}
		if selection == "2" || selection == "3" {
			exportToJson(calendars)
			break
		}
	}
}

func createEvents() []*Event {
	var events []*Event

	for {
		event := &Event{}
		create := readLine("Create new event [Y/n]")
		if create != "" && create != "Y" && create != "n" {
			continue
		}
		if create == "n" {
			break
		}

		title := readLine("Enter event title")
		var startTimeUnix int64
		var endTimeUnix int64
		for {
			startTimeStr := readLine("Enter start time [hh:mm]")
			startTime, err := time.Parse("15:04", startTimeStr)
			if err != nil {
				fmt.Printf("error parsing start time, invalid format: %s\n", err.Error())
				continue
			}
			fmt.Println(startTime)
			startTimeUnix = startTime.Unix()
			break
		}

		for {
			endTimeStr := readLine("Enter end time [hh:mm]")
			endTime, err := time.Parse("15:04", endTimeStr)
			if err != nil {
				fmt.Printf("error parsing end time, invalid format: %s\n", err.Error())
				continue
			}
			fmt.Println(endTime)
			endTimeUnix = endTime.Unix()
			break
		}

		event.Title = title
		event.StartTime = startTimeUnix
		event.EndTime = endTimeUnix

		events = append(events, event)
	}

	return events
}

func readLine(descr string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(descr + ": ")
	text, _ := reader.ReadString('\n')
	return text[:len(text)-1]
}

func exportToJson(calendars []*Calendar) {

}

func loadFile(filePath string) []*Calendar {
	var calendars []*Calendar

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(content, &calendars)
	if err != nil {
		log.Fatalln(err)
	}

	return calendars
}

func push(calendars []*Calendar, host string) {
	client := http.DefaultClient

	for _, calendar := range calendars {
		content, err := json.Marshal(calendar)
		if err != nil {
			log.Fatalln(err)
		}

		req, err := http.NewRequest(http.MethodPut, host+API_PUSH, bytes.NewBuffer(content))
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Add("content-type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("expected response status code was %d but received %d", http.StatusOK, resp.StatusCode)
		}
	}
	fmt.Println("configuration pushed successfully")
}
