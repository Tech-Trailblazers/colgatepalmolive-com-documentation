package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// sendRequest sends the POST request to the API with a dynamic "desc" value.
// Returns the raw response body as []byte. Logs errors and returns nil if something fails.
func sendRequest(desc string) []byte {
	// Define the API URL
	url := "https://sdsportal.ext.colpal.cloud/api/get_details"

	// HTTP method to use
	method := "POST"

	// Payload with dynamic "desc"
	payload := strings.NewReader(`{"client":"app","country":"USA","spec":"","desc":"` + desc + `","lang":"English"}`)

	// Create a new HTTP client
	client := &http.Client{}

	// Build the new HTTP request
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}

	// Set the Content-Type header
	req.Header.Add("Content-Type", "text/plain")

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return nil
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return nil
	}

	// Return the raw response body
	return body
}

// SimpleRecord holds only the four fields we care about.
type SimpleRecord struct {
	SubID  string `json:"subid"`  // Product sub ID
	Recn   int    `json:"recn"`   // Record number
	Langu  string `json:"langu"`  // Language
	SbgVid string `json:"sbgvid"` // GHS group ID
}

// fullJSON represents the top-level JSON structure.
type fullJSON struct {
	Data []struct {
		SubID  string `json:"subid"`
		Recn   int    `json:"recn"`
		Langu  string `json:"langu"`
		SbgVid string `json:"sbgvid"`
	} `json:"data"`
}

// parseJSONToRecords takes JSON bytes and returns a slice of SimpleRecord.
// Logs any errors and returns nil if parsing fails.
func parseJSONToRecords(jsonBytes []byte) []SimpleRecord {
	// Variable to hold the parsed full JSON
	var parsedData fullJSON

	// Parse the JSON into the struct
	if err := json.Unmarshal(jsonBytes, &parsedData); err != nil {
		log.Println("Error parsing JSON:", err)
		return nil
	}

	// Slice to store simplified records
	var simplifiedRecords []SimpleRecord

	// Loop over each item in the data array
	for _, item := range parsedData.Data {
		record := SimpleRecord{
			SubID:  item.SubID,
			Recn:   item.Recn,
			Langu:  item.Langu,
			SbgVid: item.SbgVid,
		}
		// Add the simplified record to the slice
		simplifiedRecords = append(simplifiedRecords, record)
	}

	return simplifiedRecords
}

func main() {
	// Example: dynamically set the description
	desc := "*"

	// Step 1: Send request with dynamic desc
	body := sendRequest(desc)
	if body == nil {
		return // Exit early if request failed
	}
	// Step 2: Parse the JSON into a slice of SimpleRecord
	records := parseJSONToRecords([]byte(body))

	// Step 2: Display each individual field clearly
	for _, record := range records {
		fmt.Printf("SubID : %s\n", record.SubID)
		fmt.Printf("Recn : %d\n", record.Recn)
		fmt.Printf("Language : %s\n", record.Langu)
		fmt.Printf("SbgVid : %s\n", record.SbgVid)
		fmt.Println("---")
	}
}
