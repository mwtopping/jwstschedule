package main

import (
	"context"
	"database/sql"
	"fmt"
	"jwstscheduler/internal/database"
	"log"
	"strings"
	"time"
)

type apiConfig struct {
	dbQueries *database.Queries
	db        *sql.DB
}

func (cfg *apiConfig) create_program_database() error {

	//gocy1url := "https://www.stsci.edu/jwst/science-execution/approved-programs/general-observers/cycle-1-go"
	//	urls := []string{"https://www.stsci.edu/jwst/science-execution/approved-programs/general-observers/cycle-1-go",
	//		"https://www.stsci.edu/jwst/science-execution/approved-programs/general-observers/cycle-2-go",
	//		"https://www.stsci.edu/jwst/science-execution/approved-programs/general-observers/cycle-3-go",
	//		"https://www.stsci.edu/jwst/science-execution/approved-programs/general-observers/cycle-4-go",
	//		"https://www.stsci.edu/jwst/science-execution/approved-programs/guaranteed-time-observations",
	//		"https://www.stsci.edu/jwst/science-execution/approved-programs/directors-discretionary-time"}

	urls := []string{"https://www.stsci.edu/jwst/science-execution/approved-programs/general-observers/cycle-4-go"}

	for _, url := range urls {
		log.Println("Getting programs from URL: ", url)
		programs := get_accepted_ids(url)

		for _, p := range programs {
			if p.ID != 0 {

				programParams := database.CreateProgramParams{
					ID:             int64(p.ID),
					CreatedAt:      time.Now().Unix(),
					UpdatedAt:      time.Now().Unix(),
					Title:          p.Title,
					Pi:             p.PI,
					Eap:            int64(p.EAP),
					Primetime:      float64(p.PrimeTime),
					Paralleltime:   float64(p.ParallelTime),
					Instrumentmode: strings.Join(p.InstrumentMode, ","),
					Programtype:    strings.Join(p.ProgramType, ","),
				}
				_, err := cfg.dbQueries.CreateProgram(context.Background(), programParams)
				if err != nil {
					fmt.Println("ERROR INSERTING PROGRAM TO DB: ", p.ID)
					fmt.Println(err)
				}

			}

		}
	}
	return nil
}
