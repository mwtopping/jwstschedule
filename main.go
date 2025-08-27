package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"jwstscheduler/internal/database"
	_ "modernc.org/sqlite"
)

func parse_time(t string) time.Time {
	layout := "Jan 2, 2006 15:04:05"
	parsedTime, err := time.Parse(layout, t)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return time.Unix(0, 0)
	}
	return parsedTime.UTC()
}

//go:embed sql/schema/001_program_info.sql
var program_info_schema string

//go:embed sql/schema/002_visits.sql
var visits_schema string

func main() {

	cfg := apiConfig{}

	os.Remove("./jwst.db")

	db, err := sql.Open("sqlite", "./jwst.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create tables here
	if _, err := db.ExecContext(context.Background(), program_info_schema); err != nil {
		log.Fatal(err)
	}
	if _, err := db.ExecContext(context.Background(), visits_schema); err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)
	cfg.dbQueries = queries

	err = cfg.create_program_database()
	if err != nil {
		log.Fatal("Error creating program database", err)
		return
	}

	program_IDs, err := cfg.dbQueries.GetProgramIDs(context.Background())
	if err != nil {
		log.Fatal("Error getting IDS")
		return
	}

	//cfg.get_program_info(int(8582))
	//fmt.Scanln()

	for _, ID := range program_IDs {
		cfg.get_program_info(int(ID))
		fmt.Scanln()
	}

}
