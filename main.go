package main

import (
	"encoding/json"
	"flag"
	"fmt"
	//"io/ioutil" // Needed for the json.unmarshal usage below
	"net/http"
)

// TODO: NO error handling. Should probably do something about that.
func main() {
	numResults := flag.Int("numResults", 50, "Number of results to return")
	startPage := flag.Int("startPage", 0, "Results page to start with")
	flag.Parse()

	fmt.Println("numResults: ", *numResults)
	fmt.Println("startPage: ", *startPage)

	var searchURL = fmt.Sprint("https://archive.org/advancedsearch.php?q=collection%3AGratefulDead"+
		"&fl%5B%5D=identifier&sort%5B%5D=&sort%5B%5D=&sort%5B%5D=&rows=", *numResults, 
        "&page=", *startPage, "&output=json")

	// Print search URL to the console
	fmt.Println(searchURL)

	r, _ := http.Get(searchURL)
	defer r.Body.Close()

	var searchResponse SearchResponse
	// Should this be a look until EOF? Or am I sure I have all the response?
	// Is that what the Body.Close call gets me?
	json.NewDecoder(r.Body).Decode(&searchResponse)

	// Print search result details to the console
	fmt.Println("Rows found: ", searchResponse.Response.NumFound)
	fmt.Println("Rows returned: ", searchResponse.ResponseHeader.Params.Rows)
	fmt.Println("Starting record: ", searchResponse.Response.Start)

	// Print show details to the console
	for _, doc := range searchResponse.Response.Docs {
		fmt.Println(doc.Identifier)
		var showURL = "http://archive.org/details/" + doc.Identifier + "?output=json"
		fmt.Println(showURL)
		r, _ := http.Get(showURL)
		defer r.Body.Close()

		var showResponse DeadShow
		// Should this be a look until EOF? Or am I sure I have all the response?
		// Is that what the Body.Close call gets me?
		json.NewDecoder(r.Body).Decode(&showResponse)

		// Alternative solution using Unmarshal
		// show, _ := ioutil.ReadAll(s.Body)
		// json.Unmarshal(show, &showResponse)

		// :( All these will panic if value is nil
        fmt.Println("Server: ", showResponse.Server)
		fmt.Printf("Results: %v\n", showResponse)
		fmt.Println("Date: ", showResponse.Metadata.Date[0])
		fmt.Println("Venue: ", showResponse.Metadata.Venue[0])
	}
}
