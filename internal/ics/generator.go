package ics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	ics "github.com/arran4/golang-ical"
)

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Description string `json:"description"`
}

func HandleGenerateICS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var events []Event
	err := json.NewDecoder(r.Body).Decode(&events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	icsContent, err := GenerateICS(events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", "attachment; filename=calendar.ics")
	w.Write([]byte(icsContent))
}

func GenerateLocalICS(inputFile, outputFile string) {
	jsonData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	var events []Event
	err = json.Unmarshal(jsonData, &events)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	icsContent, err := GenerateICS(events)
	if err != nil {
		log.Fatalf("Error generating ICS: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	err = ioutil.WriteFile(outputFile, []byte(icsContent), 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	log.Printf("ICS file generated successfully: %s", outputFile)
}

func GenerateICS(events []Event) (string, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)

	for _, event := range events {
		e := cal.AddEvent(fmt.Sprintf("%s@example.com", event.ID))
		e.SetCreatedTime(time.Now())
		e.SetDtStampTime(time.Now())
		e.SetModifiedAt(time.Now())
		e.SetSummary(event.Title)
		e.SetDescription(event.Description)

		start, err := time.Parse("2006-01-02", event.Start)
		if err != nil {
			return "", fmt.Errorf("invalid start date for event %s: %v", event.ID, err)
		}

		end, err := time.Parse("2006-01-02", event.End)
		if err != nil {
			return "", fmt.Errorf("invalid end date for event %s: %v", event.ID, err)
		}
		end = end.AddDate(0, 0, 1) // Add one day to make it inclusive

		e.SetAllDayStartAt(start)
		e.SetAllDayEndAt(end)
	}

	return cal.Serialize(), nil
}
