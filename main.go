package main

import (
	"fmt"
	"io"
	//"net/http"
	"encoding/xml"
	"os"
	"time"
)

func parse_time(t string) time.Time {
	layout := "Jan 2, 2006 15:04:05"
	parsedTime, err := time.Parse(layout, t)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return time.Time{}
	}
	return parsedTime.UTC()
}

func main() {
	fmt.Println("Hello")
	programs := get_accepted_ids()

	for _, p := range programs {
		p.Show()
	}

	//	PARSE_URL := "https://www.stsci.edu/jwst-program-info/visits/?program=1914"
	PARSE_URL := "https://www.stsci.edu/jwst-program-info/visits/?program=1914&download=&pi=1&referrer=https://www.stsci.edu"
	fmt.Println("Parsing ", PARSE_URL)

	response, err := os.Open("./test.xml")
	//response, err := http.Get(PARSE_URL)
	if err != nil {
		fmt.Println("Error Receiving HTML")
		return
	}
	//defer response.Body.Close()
	defer response.Close()
	//body, err := io.ReadAll(response.Body)
	body, err := io.ReadAll(response)

	program_info := VisitStatusReport{}

	err = xml.Unmarshal(body, &program_info)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, visit := range program_info.Visit {
		fmt.Printf("Obs:%v, Vis:%v, Status:%v, StartTime:%v\n",
			visit.Observation,
			visit.Visit,
			visit.Status,
			parse_time(visit.StartTime).Unix())
	}

}
