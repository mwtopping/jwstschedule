package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"jwstscheduler/internal/database"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (cfg *apiConfig) get_program_info(ID int) {
	url := get_program_url(ID)
	fmt.Println("Parsing ", url)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error Receiving HTML")
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	program_info := VisitStatusReport{}

	err = xml.Unmarshal(body, &program_info)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, visit := range program_info.Visit {
		//fmt.Printf("Obs:%v, Vis:%v, Status:%v, StartTime:%v, EndTime:%v\n",
		//	visit.Observation,
		//	visit.Visit,
		//	visit.Status,
		//	visit.StartTime,
		//	visit.EndTime)
		//fmt.Println(visit)

		cfg.add_visit(ID, visit)
	}
}

func (cfg *apiConfig) add_visit(programID int, visit SingleVisit) error {

	obsInt, err := strconv.Atoi(visit.Observation)
	if err != nil {
		return err
	}
	visInt, err := strconv.Atoi(visit.Visit)
	if err != nil {
		return err
	}

	// Format each number with appropriate zero-padding
	formatted := fmt.Sprintf("%d%03d%03d", programID, obsInt, visInt)

	// Convert back to integer
	fullID, err := strconv.Atoi(formatted)
	if err != nil {
		fmt.Println("Error converting to integer:", err)
		return err
	}

	var sTime int64
	var eTime int64

	switch visit.Status {
	case "Flight Ready", "Implementation":
		// need to parse planwindow into a start and end time
		timeRange := visit.PlanWindow
		timeRange, _ = get_substring_between(timeRange, "(", ")")
		endPoints := strings.Split(timeRange, "-")
		if len(endPoints) != 2 {
			sTime = 0
			eTime = 0
		} else {
			sTimeYears := strings.TrimSpace(endPoints[0])
			//			sTimeYears, err := strconv.ParseFloat(strings.TrimSpace(endPoints[0]), 64)
			//			if err != nil {
			//				fmt.Println(err)
			//				return err
			//			}
			eTimeYears := strings.TrimSpace(endPoints[1])
			//			eTimeYears, err := strconv.ParseFloat(strings.TrimSpace(endPoints[1]), 64)
			//			if err != nil {
			//				fmt.Println(err)
			//				return err
			//			}

			sTime = fractionalYearToUnix(sTimeYears)
			eTime = fractionalYearToUnix(eTimeYears)
		}
	case "Withdrawn":
		sTime = 0
		eTime = 0
	default:
		sTime = parse_time(visit.StartTime).Unix()
		eTime = parse_time(visit.EndTime).Unix()
	}

	visit_params := database.CreateVisitParams{
		ID:            int64(fullID),
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
		ProgramID:     int64(programID),
		Observation:   int64(obsInt),
		Visit:         int64(visInt),
		Status:        visit.Status,
		Target:        visit.Target,
		Configuration: visit.Configuration,
		Starttime:     sTime,
		Endtime:       eTime,
	}

	_, err = cfg.dbQueries.CreateVisit(context.Background(), visit_params)

	return err
}

func get_program_url(ID int) string {

	url := fmt.Sprintf("https://www.stsci.edu/jwst-program-info/visits/?program=%d&download=&pi=1&referrer=https://www.stsci.edu", ID)

	return url
}

func fractionalYearToUnix(fractionalYear string) int64 {

	date_parts := strings.Split(fractionalYear, ".")

	if len(date_parts) != 2 {
		log.Println("Not enough parts of the date")
		return 0
	}

	year, err := strconv.Atoi(date_parts[0])
	if err != nil {
		log.Println("Error converting year to int", err)
		return 0
	}

	day, err := strconv.Atoi(date_parts[1])
	if err != nil {
		log.Println("Error converting year to int", err)
		return 0
	}

	startOfYear := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

	resultTime := startOfYear.AddDate(0, 0, day)

	return resultTime.Unix()
}
