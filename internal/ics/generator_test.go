package ics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGenerateICSFromPublicHolidays(t *testing.T) {
	// Read test data
	testDataPath, _ := filepath.Abs("../../testdata/public_holidays_2024.json")
	jsonFile, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var events []Event
	err = json.Unmarshal(jsonFile, &events)
	if err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	// Generate ICS
	icsContent, err := GenerateICS(events)
	if err != nil {
		t.Fatalf("Error generating ICS: %v", err)
	}

	// Print ICS content for debugging
	t.Logf("Generated ICS content:\n%s", icsContent)

	// Basic checks
	if !strings.Contains(icsContent, "BEGIN:VCALENDAR") {
		t.Error("ICS content does not contain BEGIN:VCALENDAR")
	}
	if !strings.Contains(icsContent, "END:VCALENDAR") {
		t.Error("ICS content does not contain END:VCALENDAR")
	}

	// Check for each holiday
	for _, event := range events {
		if !strings.Contains(icsContent, "SUMMARY:"+event.Title) {
			t.Errorf("Event '%s' not found in ICS content", event.Title)
		}
		expectedStart := "DTSTART;VALUE=DATE:" + strings.ReplaceAll(event.Start, "-", "")
		if !strings.Contains(icsContent, expectedStart) {
			t.Errorf("Start date for event '%s' not found or incorrect in ICS content. Expected %s", event.Title, expectedStart)
		}
		endDate, _ := time.Parse("2006-01-02", event.End)
		endDate = endDate.AddDate(0, 0, 1) // Add one day because ICS end date is exclusive
		expectedEnd := "DTEND;VALUE=DATE:" + endDate.Format("20060102")
		if !strings.Contains(icsContent, expectedEnd) {
			t.Errorf("End date for event '%s' not found or incorrect in ICS content. Expected %s", event.Title, expectedEnd)
		}
	}
}

func TestGenerateICS(t *testing.T) {
	// Test data
	events := []Event{
		{
			ID:          "event1",
			Title:       "Test Event",
			Start:       "2023-08-01",
			End:         "2023-08-02",
			Description: "This is a test event",
		},
	}

	// Generate ICS
	icsContent, err := GenerateICS(events)
	if err != nil {
		t.Fatalf("Error generating ICS: %v", err)
	}

	// Basic checks
	if !strings.Contains(icsContent, "BEGIN:VCALENDAR") {
		t.Error("ICS content does not contain BEGIN:VCALENDAR")
	}
	if !strings.Contains(icsContent, "END:VCALENDAR") {
		t.Error("ICS content does not contain END:VCALENDAR")
	}
	if !strings.Contains(icsContent, "BEGIN:VEVENT") {
		t.Error("ICS content does not contain BEGIN:VEVENT")
	}
	if !strings.Contains(icsContent, "END:VEVENT") {
		t.Error("ICS content does not contain END:VEVENT")
	}

	// Check event details
	if !strings.Contains(icsContent, "SUMMARY:Test Event") {
		t.Error("Event summary not found in ICS content")
	}
	if !strings.Contains(icsContent, "DESCRIPTION:This is a test event") {
		t.Error("Event description not found in ICS content")
	}
	if !strings.Contains(icsContent, "DTSTART;VALUE=DATE:20230801") {
		t.Error("Event start date not found or incorrect in ICS content")
	}
	if !strings.Contains(icsContent, "DTEND;VALUE=DATE:20230803") {
		t.Error("Event end date not found or incorrect in ICS content")
	}

	// Check for required fields
	requiredFields := []string{"UID", "DTSTAMP", "DTSTART", "DTEND", "SUMMARY"}
	for _, field := range requiredFields {
		if !strings.Contains(icsContent, field) {
			t.Errorf("Required field %s not found in ICS content", field)
		}
	}

	// Validate date format
	dateStr := time.Now().Format("20060102T150405Z")
	if !strings.Contains(icsContent, "DTSTAMP:"+dateStr[:8]) {
		t.Error("DTSTAMP date format is incorrect")
	}
}
