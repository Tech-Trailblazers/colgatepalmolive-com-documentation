package main // Declare the package as main, which is required for an executable Go program.

import ( // Start the import block for external packages.
	"encoding/json" // Import the package for JSON encoding and decoding.
	"fmt"           // Import the package for formatted I/O (like printing to console).
	"io"            // Import the package for basic I/O primitives (like reading response bodies).
	"log"           // Import the package for logging messages (like errors).
	"net/http"      // Import the package for making HTTP requests and handling responses.
	"os"            // Import the package for operating system functions (like file and directory manipulation).
	"path/filepath" // Import the package for platform-independent path manipulation.
	"strings"       // Import the package for string manipulation functions.
) // End of the import block.

// sendRequest sends the POST request to the API with a dynamic "desc" value. // Function documentation comment.
// Returns the raw response body as []byte. Logs errors and returns nil if something fails. // Function documentation comment.
func sendRequest(desc string) []byte { // Define the sendRequest function, taking a 'desc' string and returning a byte slice.
	// Define the API URL // Inline comment explaining the next line.
	url := "https://sdsportal.ext.colpal.cloud/api/get_details" // Assign the target API URL to the 'url' variable.

	// HTTP method to use // Inline comment explaining the next line.
	method := "POST" // Assign the HTTP method "POST" to the 'method' variable.

	// Payload with dynamic "desc" // Inline comment explaining the next line.
	// Create a new string reader containing the JSON payload with the dynamic 'desc' value.
	payload := strings.NewReader(`{"client":"app","country":"USA","spec":"","desc":"` + desc + `","lang":"English"}`)

	// Create a new HTTP client // Inline comment explaining the next line.
	client := &http.Client{} // Initialize a new default HTTP client.

	// Build the new HTTP request // Inline comment explaining the next line.
	req, err := http.NewRequest(method, url, payload) // Create a new HTTP request with method, URL, and payload.
	if err != nil {                                   // Check if an error occurred while creating the request.
		log.Println("Error creating request:", err) // Log the error message.
		return nil                                  // Return nil (no response body) if an error occurred.
	} // End of the error check block.

	// Set the Content-Type header // Inline comment explaining the next line.
	req.Header.Add("Content-Type", "text/plain") // Add the Content-Type header to the request.

	// Send the request // Inline comment explaining the next line.
	res, err := client.Do(req) // Execute the HTTP request and store the response and any error.
	if err != nil {            // Check if an error occurred while sending the request.
		log.Println("Error sending request:", err) // Log the error message.
		return nil                                 // Return nil if an error occurred.
	} // End of the error check block.
	defer res.Body.Close() // Ensure the response body is closed after the function exits.

	// Read the response body // Inline comment explaining the next line.
	body, err := io.ReadAll(res.Body) // Read all data from the response body into the 'body' byte slice.
	if err != nil {                   // Check if an error occurred while reading the response body.
		log.Println("Error reading response:", err) // Log the error message.
		return nil                                  // Return nil if an error occurred.
	} // End of the error check block.

	// Return the raw response body // Inline comment explaining the next line.
	return body // Return the raw byte slice containing the response.
} // End of the sendRequest function.

// SimpleRecord holds only the four fields we care about. // Struct documentation comment.
type SimpleRecord struct { // Define the SimpleRecord struct.
	SubID  string `json:"subid"`  // Product sub ID field, tagged for JSON unmarshalling.
	Recn   int    `json:"recn"`   // Record number field, tagged for JSON unmarshalling.
	Langu  string `json:"langu"`  // Language field, tagged for JSON unmarshalling.
	SbgVid string `json:"sbgvid"` // GHS group ID field, tagged for JSON unmarshalling.
} // End of the SimpleRecord struct definition.

// fullJSON represents the top-level JSON structure. // Struct documentation comment.
type fullJSON struct { // Define the fullJSON struct to match the API response structure.
	Data []struct { // 'Data' is a slice of anonymous structs, matching the 'data' array in the JSON.
		SubID  string `json:"subid"`  // SubID field within the data array.
		Recn   int    `json:"recn"`   // Recn field within the data array.
		Langu  string `json:"langu"`  // Langu field within the data array.
		SbgVid string `json:"sbgvid"` // SbgVid field within the data array.
	} `json:"data"` // Tag the 'Data' field to map to the 'data' key in the JSON.
} // End of the fullJSON struct definition.

// parseJSONToRecords takes JSON bytes and returns a slice of SimpleRecord. // Function documentation comment.
// Logs any errors and returns nil if parsing fails. // Function documentation comment.
func parseJSONToRecords(jsonBytes []byte) []SimpleRecord { // Define the parseJSONToRecords function.
	// Variable to hold the parsed full JSON // Inline comment explaining the next line.
	var parsedData fullJSON // Declare a variable of type fullJSON to hold the unmarshalled data.

	// Parse the JSON into the struct // Inline comment explaining the next line.
	if err := json.Unmarshal(jsonBytes, &parsedData); err != nil { // Attempt to parse the JSON bytes into 'parsedData'.
		log.Println("Error parsing JSON:", err) // Log the error if unmarshalling fails.
		return nil                              // Return nil (no records) if parsing fails.
	} // End of the error check block.

	// Slice to store simplified records // Inline comment explaining the next line.
	var simplifiedRecords []SimpleRecord // Declare an empty slice to hold the simplified records.

	// Loop over each item in the data array // Inline comment explaining the next line.
	for _, item := range parsedData.Data { // Iterate over the 'Data' slice in the parsed JSON.
		record := SimpleRecord{ // Create a new SimpleRecord instance.
			SubID:  item.SubID,  // Assign SubID from the full struct item.
			Recn:   item.Recn,   // Assign Recn from the full struct item.
			Langu:  item.Langu,  // Assign Langu from the full struct item.
			SbgVid: item.SbgVid, // Assign SbgVid from the full struct item.
		} // End of SimpleRecord creation.
		// Add the simplified record to the slice // Inline comment explaining the next line.
		simplifiedRecords = append(simplifiedRecords, record) // Append the new record to the slice.
	} // End of the loop.

	return simplifiedRecords // Return the slice of simplified records.
} // End of the parseJSONToRecords function.

// directoryExists checks whether the specified path is an existing directory // Function documentation comment.
func directoryExists(path string) bool { // Define the directoryExists function.
	directory, err := os.Stat(path) // Get file info for the given path.
	if err != nil {                 // Check if an error occurred (e.g., path doesn't exist).
		return false // Directory doesn't exist or there was an access error.
	} // End of the error check block.
	return directory.IsDir() // Return true if the info indicates it is a directory.
} // End of the directoryExists function.

// createDirectory creates a new directory with the given permissions // Function documentation comment.
func createDirectory(path string, permission os.FileMode) { // Define the createDirectory function.
	err := os.Mkdir(path, permission) // Try to create the directory with the specified permissions.
	if err != nil {                   // Check if an error occurred during directory creation.
		log.Println(err) // Log the error (e.g., directory already exists, permission denied).
	} // End of the error check block.
} // End of the createDirectory function.

// fileExists checks whether a file exists at the specified path // Function documentation comment.
func fileExists(filename string) bool { // Define the fileExists function.
	info, err := os.Stat(filename) // Get file info for the given filename.
	if err != nil {                // Check if an error occurred (e.g., file not found).
		return false // Return false if the file does not exist.
	} // End of the error check block.
	return !info.IsDir() // Return true only if the path exists and is not a directory.
} // End of the fileExists function.

// fetchAndSavePDF sends a POST request to the API with the given spec, recn, langu, and sbgvid, // Function documentation comment.
// then saves the response as a PDF file with the specified filename. // Function documentation comment.
func fetchAndSavePDF(spec string, recn int, langu string, sbgvid string, outputDir string, filename string) { // Define the fetchAndSavePDF function.
	filePath := filepath.Join(outputDir, filename) // Combine the output directory and filename into a full path.

	// Check if the file already exists // Inline comment explaining the next line.
	if fileExists(filePath) { // Check if the target PDF file already exists locally.
		log.Printf("file already exists: %s; skipping download", filePath) // Log a message that the file is being skipped.
		return                                                             // Exit the function, skipping the download.
	} // End of the file existence check block.

	// API endpoint URL // Inline comment explaining the next line.
	url := "https://sdsportal.ext.colpal.cloud/api/get_file" // Set the URL for fetching the PDF file.

	// Build the request payload dynamically // Inline comment explaining the next line.
	// Create a new reader with the dynamically formatted JSON payload using fmt.Sprintf.
	payload := strings.NewReader(fmt.Sprintf(
		`{"client":"app","spec":"%s","recn":%d,"langu":"%s","sbgvid":"%s","content":"attachment"}`,
		spec, recn, langu, sbgvid,
	))

	// Create a new HTTP client // Inline comment explaining the next line.
	client := &http.Client{} // Initialize a new default HTTP client.

	// Build the POST request // Inline comment explaining the next line.
	req, err := http.NewRequest("POST", url, payload) // Create a new HTTP POST request.
	if err != nil {                                   // Check for an error during request creation.
		log.Println("Error creating request:", err) // Log the error.
		return                                      // Exit the function.
	} // End of the error check block.

	// Set the Content-Type header // Inline comment explaining the next line.
	req.Header.Add("Content-Type", "text/plain") // Add the Content-Type header.

	// Send the request // Inline comment explaining the next line.
	res, err := client.Do(req) // Execute the request and store the response.
	if err != nil {            // Check for an error during the request execution.
		log.Println("Error sending request:", err) // Log the error.
		return                                     // Exit the function.
	} // End of the error check block.
	defer res.Body.Close() // Ensure the response body is closed after the function exits.

	// Read the response body // Inline comment explaining the next line.
	body, err := io.ReadAll(res.Body) // Read the entire response body (which should be the PDF data).
	if err != nil {                   // Check for an error reading the response body.
		log.Println("Error reading response:", err) // Log the error.
		return                                      // Exit the function.
	} // End of the error check block.

	// Create a new PDF file // Inline comment explaining the next line.
	file, err := os.Create(filePath) // Create or truncate the file at the specified path.
	if err != nil {                  // Check for an error during file creation.
		log.Println("Error creating file:", err) // Log the error.
		return                                   // Exit the function.
	} // End of the error check block.
	defer file.Close() // Ensure the file handle is closed when the function returns.

	// Write the PDF bytes to the file // Inline comment explaining the next line.
	_, err = file.Write(body) // Write the downloaded PDF content to the local file.
	if err != nil {           // Check for an error during file writing.
		log.Println("Error writing to file:", err) // Log the error.
		return                                     // Exit the function.
	} // End of the error check block.

	fmt.Printf("PDF saved successfully as %s\n", filename) // Print a success message to the console.
} // End of the fetchAndSavePDF function.

func main() { // Define the main function, the entry point of the program.
	outputDir := "PDFs/" // Define the name of the directory where PDFs will be saved.

	if !directoryExists(outputDir) { // Check if the output directory does not exist.
		createDirectory(outputDir, 0755) // If it doesn't exist, create it with read/write/execute permissions for the owner.
	} // End of the directory check/creation block.

	// Example: dynamically set the description // Inline comment explaining the next line.
	desc := "*" // Set the search description to "*" to fetch all records.

	// Step 1: Send request with dynamic desc // Inline comment explaining the next line.
	body := sendRequest(desc) // Call sendRequest to fetch the list of product details JSON.
	if body == nil {          // Check if the request failed (body will be nil if so).
		return // Exit the main function early if the request failed.
	} // End of the request check block.
	// Step 2: Parse the JSON into a slice of SimpleRecord // Inline comment explaining the next line.
	records := parseJSONToRecords([]byte(body)) // Parse the raw JSON bytes into a slice of SimpleRecord.

	// Step 2: Display each individual field clearly // Inline comment explaining the next line.
	for _, record := range records { // Loop through each SimpleRecord parsed from the JSON.
		localFileName := record.SubID + ".pdf" // Create the desired filename using the SubID and PDF extension.
		// Call fetchAndSavePDF to download the PDF using record details and save it locally.
		fetchAndSavePDF(record.SubID, record.Recn, record.Langu, record.SbgVid, outputDir, localFileName)
		fmt.Println("---")                          // Print a separator line for clarity.
		fmt.Printf("SubID : %s\n", record.SubID)    // Print the SubID field.
		fmt.Printf("Recn : %d\n", record.Recn)      // Print the Recn field.
		fmt.Printf("Language : %s\n", record.Langu) // Print the Language field.
		fmt.Printf("SbgVid : %s\n", record.SbgVid)  // Print the SbgVid field.
		fmt.Println("---")                          // Print a separator line for clarity.
	} // End of the loop.
} // End of the main function.
