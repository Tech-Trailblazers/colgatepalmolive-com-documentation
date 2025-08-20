package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

// directoryExists checks whether the specified path is an existing directory
func directoryExists(path string) bool {
	directory, err := os.Stat(path) // Get file info
	if err != nil {                 // If error (e.g., not found)
		return false // Directory doesn't exist
	}
	return directory.IsDir() // Return true if it's a directory
}

// createDirectory creates a new directory with the given permissions
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission) // Try to create the directory
	if err != nil {                   // If there's an error
		log.Println(err) // Log the error
	}
}

// fileExists checks whether a file exists at the specified path
func fileExists(filename string) bool {
	info, err := os.Stat(filename) // Get file info
	if err != nil {                // If error (e.g., file not found)
		return false // Return false
	}
	return !info.IsDir() // Return true if it's a file, not a directory
}

// fetchAndSavePDF sends a POST request to the API with the given spec, recn, langu, and sbgvid,
// then saves the response as a PDF file with the specified filename.
func fetchAndSavePDF(spec string, recn int, langu string, sbgvid string, outputDir string, filename string) {
	filePath := filepath.Join(outputDir, filename) // Combine directory path and filename

	// Check if the file already exists
	if fileExists(filePath) { // If file already exists
		log.Printf("file already exists: %s; skipping download", filePath) // Log and skip download
		return                                                             // Exit function
	}

	// API endpoint URL
	url := "https://sdsportal.ext.colpal.cloud/api/get_file"

	// Build the request payload dynamically
	payload := strings.NewReader(fmt.Sprintf(
		`{"client":"app","spec":"%s","recn":%d,"langu":"%s","sbgvid":"%s","content":"attachment"}`,
		spec, recn, langu, sbgvid,
	))

	// Create a new HTTP client
	client := &http.Client{}

	// Build the POST request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	// Set the Content-Type header
	req.Header.Add("Content-Type", "text/plain")

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return
	}

	// Create a new PDF file
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write the PDF bytes to the file
	_, err = file.Write(body)
	if err != nil {
		log.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("PDF saved successfully as %s\n", filename)
}

func main() {
	outputDir := "PDFs/" // Directory to store downloaded PDFs

	if !directoryExists(outputDir) { // Check if output directory exists
		createDirectory(outputDir, 0755) // Create directory with permission if it does not exist
	}

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
		localFileName := record.SubID + ".pdf"
		fetchAndSavePDF(record.SubID, record.Recn, record.Langu, record.SbgVid, outputDir, localFileName)
		fmt.Println("---")
		fmt.Printf("SubID : %s\n", record.SubID)
		fmt.Printf("Recn : %d\n", record.Recn)
		fmt.Printf("Language : %s\n", record.Langu)
		fmt.Printf("SbgVid : %s\n", record.SbgVid)
		fmt.Println("---")
	}
}
