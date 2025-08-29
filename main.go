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

	err := cfg.ResetDatabase("./jwst.db")
	if err != nil {
		log.Fatal(err)
	}
	defer cfg.db.Close()

	//err := cfg.LoadDatabase("./jwst.db")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer cfg.db.Close()

	program_IDs, err := cfg.dbQueries.GetProgramIDs(context.Background())
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error getting IDS")
	}

	for _, p := range program_IDs {
		fmt.Println(p)
	}

}

func (cfg *apiConfig) ResetDatabase(s string) error {
	os.Remove(s)
	db, err := sql.Open("sqlite", s)
	cfg.db = db
	if err != nil {
		log.Fatal(err)
	}

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
		return err
	}

	program_IDs, err := cfg.dbQueries.GetProgramIDs(context.Background())
	if err != nil {
		log.Fatal("Error getting IDS")
		return err
	}

	interval := time.Duration(10) * time.Second
	ticker := time.NewTicker(interval)

	for _, ID := range program_IDs {
		fmt.Println(ID)
		<-ticker.C
		cfg.get_program_info(int(ID))
		//fmt.Scanln()
	}
	return nil
}

func (cfg *apiConfig) LoadDatabase(s string) error {
	db, err := sql.Open("sqlite", s)
	cfg.db = db
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)
	cfg.dbQueries = queries

	return nil
}
