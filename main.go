package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Handler for the home page.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hi Snippetbox"))
}

// Handler for displaying the content of a note.
func showSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the parameter id from the URL and try
	// to convert the string to an integer using the strconv.Atoi() function. If it cannot
	// be converted to an integer, or if the value is less than 1, return a response
	// 404 - page not found!
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Fprintf() function to insert the value from id into the response string
	// and write it to http.ResponseWriter.
	fmt.Fprintf(w, "Displaying the selected note with ID %d...", id)
}

// Handler for creating a new note.
func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check if the request uses the POST method
	if r.Method != http.MethodPost {
		// Use the Header().Set() method to add the 'Allow: POST' header to
		// the HTTP header map. The first parameter is the name of the header, and
		// the second parameter is the value of the header.
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "GET Method forbidden!\n", 405)
		return
	}
	w.Write([]byte("Add new Snippet"))
}

func main() {
	// Register two new handlers and their corresponding URL patterns in
	// the servemux router
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Running the web server on http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
