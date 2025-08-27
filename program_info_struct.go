package main

import (
	"encoding/xml"
)

type VisitStatusReport struct {
	XMLName     xml.Name `xml:"visitStatusReport"`
	Text        string   `xml:",chardata"`
	Observatory string   `xml:"observatory,attr"`
	ID          string   `xml:"id,attr"`
	Title       string   `xml:"title"`
	ReportTime  string   `xml:"reportTime"`
	Visit       []struct {
		Text          string `xml:",chardata"`
		Observation   string `xml:"observation,attr"`
		Visit         string `xml:"visit,attr"`
		Status        string `xml:"status"`
		Target        string `xml:"target"`
		Configuration string `xml:"configuration"`
		Hours         string `xml:"hours"`
		StartTime     string `xml:"startTime"`
		EndTime       string `xml:"endTime"`
	} `xml:"visit"`
}
