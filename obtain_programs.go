package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	//	"net/http"
	"os"
	"regexp"
)

type ProgramInfo struct {
	ID             int
	Title          string
	PI             string
	EAP            int
	PrimeTime      float32
	ParallelTime   float32
	InstrumentMode []string
	ProgramType    []string
}

func (p *ProgramInfo) Show() {
	fmt.Printf("ID:             %v\n", p.ID)
	fmt.Printf("Title:          %v\n", p.Title)
	fmt.Printf("PI:             %v\n", p.PI)
	fmt.Printf("EAP:            %v\n", p.EAP)
	fmt.Printf("PrimeTime:      %v\n", p.PrimeTime)
	fmt.Printf("ParallelTime:   %v\n", p.ParallelTime)
	fmt.Printf("InstrumentMode: %v\n", p.InstrumentMode)
	fmt.Printf("ProgramType:    %v\n", p.ProgramType)
}

func get_PI(s string) string {
	pI_index := strings.Index(s, "PI")
	if pI_index == -1 {
		return ""
	}

	// do some random cleaning
	s = strings.Replace(s, "PI:", "", -1)
	s = strings.Replace(s, "<br />", "", -1)

	coI_index := strings.Index(s, "Co-")
	if coI_index == -1 {
		return strings.TrimSpace(s)
	} else {
		return strings.TrimSpace(s[:coI_index])
	}

}

func parse_instrument_mode(s string) []string {

	instRegex := regexp.MustCompile(`<br />`)
	cleanContent := instRegex.ReplaceAllString(s, "")
	cleanContent = strings.TrimSpace(cleanContent)

	deSpaced := strings.Join(strings.Fields(strings.TrimSpace(cleanContent)), " ")

	return strings.Fields(deSpaced)
}

func parse_exptime(s string) (float32, float32) {
	split_index := strings.Index(s, "/")

	// only primetime
	if split_index == -1 {
		priT, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return float32(0.0), float32(0.0)
		}

		return float32(priT), float32(0.0)
	} else { // both prime and parallel time
		//for pure parallel, primetime should still be listed as zero

		priT, err := strconv.ParseFloat(s[:split_index], 32)
		if err != nil {
			priT = 0.0
		}

		parT, err := strconv.ParseFloat(s[split_index+1:], 32)
		if err != nil {
			parT = 0.0
		}
		return float32(priT), float32(parT)
	}
}

func parse_program_type(s string) []string {

	split_index := strings.Index(s, ",")
	if split_index == -1 {
		return []string{s}
	}
	var modes []string

	for split_index != -1 {
		modes = append(modes, strings.TrimSpace(s[:split_index]))
		s = s[split_index+1:]
		split_index = strings.Index(s, ",")
	}

	modes = append(modes, strings.TrimSpace(s))
	return modes
}

func assemble_programinfo(entries []string) ProgramInfo {
	if len(entries) != 7 {
		return ProgramInfo{}
	}

	ID, err := strconv.Atoi(entries[0])
	if err != nil {
		fmt.Println("ERROR converting ID to integer", err)
		return ProgramInfo{}
	}

	EAP, err := strconv.Atoi(entries[3])
	if err != nil {
		fmt.Println("ERROR converting EAP to integer", err)
		EAP = -1
	}

	priT, parT := parse_exptime(entries[4])

	instmodes := parse_instrument_mode(entries[5])

	progtypes := parse_program_type(entries[6])

	return ProgramInfo{ID: ID,
		Title:          entries[1],
		PI:             get_PI(entries[2]),
		EAP:            EAP,
		PrimeTime:      priT,
		ParallelTime:   parT,
		InstrumentMode: instmodes,
		ProgramType:    progtypes,
	}

}

func parse_table_row(htmlRow string) []string {
	var entries []string

	newlineRegex := regexp.MustCompile(`\n`)
	htmlRow = newlineRegex.ReplaceAllString(htmlRow, "")
	htmlRow = strings.TrimSpace(htmlRow)

	// get all info between td tags
	tdRegex := regexp.MustCompile(`<td[^>]*>(.*?)</td>`)
	matches := tdRegex.FindAllStringSubmatch(htmlRow, -1)

	for _, match := range matches {
		if len(match) > 1 {
			content := match[1]

			tagRegex := regexp.MustCompile(`<[^>]*>`)
			cleanContent := tagRegex.ReplaceAllString(content, "")
			cleanContent = strings.TrimSpace(cleanContent)

			entries = append(entries, cleanContent)
		}
	}

	return entries
}

func get_accepted_ids() []ProgramInfo {
	//PARSE_URL := "https://www.stsci.edu/jwst/science-execution/approved-programs/general-observers/cycle-1-go"

	response, err := os.Open("./programs.txt")
	//response, err := http.Get(PARSE_URL)
	if err != nil {
		fmt.Println("Error Receiving HTML")
		return []ProgramInfo{}
	}
	//defer response.Body.Close()
	defer response.Close()
	//body, err := io.ReadAll(response.Body)
	body, err := io.ReadAll(response)

	var all_programs []ProgramInfo

	body_string := string(body)
	// find each accordion in the body
	for len(body_string) > 0 {

		accordion_content, accordion_length := get_substring_between(body_string, "\"accordion-header\"", "\"accordion-header\"")

		// this gets the header
		_, n := get_substring_between(accordion_content, "<tr>", "</tr>")

		// get content of the table
		table_body, n := get_substring_between(accordion_content[n:], "<tbody>", "</tbody>")
		nrows := 0
		for len(table_body) > 0 {
			// this will get each row
			row, n := get_substring_between(table_body, "<tr>", "</tr>")
			entries := parse_table_row(row)
			program_info := assemble_programinfo(entries)
			all_programs = append(all_programs, program_info)
			table_body = table_body[n:]
			nrows += 1
			//			fmt.Scanln()
		}

		body_string = body_string[accordion_length:]
	}

	return all_programs

}

func get_substring_between(full_string, sub1, sub2 string) (string, int) {
	start_ind := strings.Index(full_string, sub1)
	if start_ind == -1 {
		return full_string, len(full_string)
	}
	end_ind := strings.Index(full_string[start_ind+len(sub1):], sub2)
	if end_ind == -1 {
		return full_string[start_ind+len(sub1):], len(full_string[start_ind+len(sub1):])
	}
	substring := full_string[start_ind+len(sub1) : start_ind+len(sub1)+end_ind]
	return substring, start_ind + len(sub1) + end_ind
}
