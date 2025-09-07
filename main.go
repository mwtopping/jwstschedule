package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
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

	const port = ":8080"

	cfg := apiConfig{}

	//	err := cfg.ResetDatabase("./jwst.db")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer cfg.db.Close()

	err := cfg.LoadDatabase("./jwst.db")
	if err != nil {
		log.Fatal(err)
	}
	defer cfg.db.Close()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", cfg.handlerDisplay)

	server := &http.Server{Handler: mux, Addr: port}
	log.Printf("Server listening on port%v\n", port)
	log.Fatal(server.ListenAndServe())

}

func (cfg *apiConfig) handlerDisplay(w http.ResponseWriter, r *http.Request) {

	//load html template
	tmpl, err := template.ParseFiles("./templates/mytemplate.html")
	if err != nil {
		log.Println("Error reading template")
		log.Println(err)
		return
	}

	type Visit struct {
		ProgID    int
		ObsNum    int
		VisNum    int
		EAP       int
		ProgName  string
		Status    string
		Starttime string
		Endtime   string
	}

	DisplayVisits := make([]Visit, 0, 0)

	var vs []database.GetAllVisitsRow

	requestPath := strings.TrimPrefix(r.URL.Path, "/")
	switch requestPath {
	case "week":
		weekvs, err := cfg.dbQueries.GetWeekVisits(context.Background(), time.Now().Unix())
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, weekv := range weekvs {
			vs = append(vs, database.GetAllVisitsRow(weekv))
		}

	default:
		vs, err = cfg.dbQueries.GetAllVisits(context.Background())

	}

	for _, v := range vs {
		startTime := time.Unix(v.Starttime, 0)
		endTime := time.Unix(v.Endtime, 0)

		DisplayVisit := Visit{
			ProgID:    int(v.ID),
			ObsNum:    int(v.Observation),
			VisNum:    int(v.Visit),
			EAP:       int(v.Eap),
			ProgName:  v.Title,
			Status:    v.Status,
			Starttime: startTime.In(time.UTC).Format("2006-01-02T15:04:05"),
			Endtime:   endTime.In(time.UTC).Format("2006-01-02T15:04:05"),
		}

		fmt.Println(startTime, v.Starttime)
		fmt.Println(startTime.In(time.UTC).Format("2006-01-02T15:04:05 MST"))

		DisplayVisits = append(DisplayVisits, DisplayVisit)

	}

	type PageData struct {
		DisplayVisits []Visit
	}

	data := PageData{DisplayVisits: DisplayVisits}

	tmpl.Execute(w, data)

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

	interval := time.Duration(4) * time.Second
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
