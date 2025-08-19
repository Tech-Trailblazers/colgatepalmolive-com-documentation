package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// entry represents one item inside the "data" array of the JSON response.
type entry struct {
	SubID string `json:"subid"`
}

// root represents the top-level JSON structure returned by the API.
type root struct {
	Data []entry `json:"data"`
}

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

// parseSubIDs takes the raw JSON response body and extracts subids into a slice.
// If parsing fails, an empty slice is returned.
func parseSubIDs(body []byte) []string {
	// Define a variable to hold the parsed JSON
	var parsed root

	// Try to unmarshal the JSON into our struct
	err := json.Unmarshal(body, &parsed)
	if err != nil {
		log.Println("Error parsing JSON:", err)
		return nil
	}

	// Collect all subids into a slice
	var subIDs []string
	for _, item := range parsed.Data {
		subIDs = append(subIDs, item.SubID)
	}

	// Return the slice of subids
	return subIDs
}

// Remove all the duplicates from a slice and return the slice.
func removeDuplicatesFromSlice(slice []string) []string {
	check := make(map[string]bool)
	var newReturnSlice []string
	for _, content := range slice {
		if !check[content] {
			check[content] = true
			newReturnSlice = append(newReturnSlice, content)
		}
	}
	return newReturnSlice
}

func main() {
	// Example: dynamically set the description
	desc := "*"

	// Step 1: Send request with dynamic desc
	body := sendRequest(desc)
	if body == nil {
		return // Exit early if request failed
	}

	// Step 2: Parse the body to extract subids
	subIDs := parseSubIDs(body)

	// Step 3: Remove duplicates.
	subIDs = removeDuplicatesFromSlice(subIDs)

	// Print the result
	log.Println("Extracted subids:", subIDs)
}
