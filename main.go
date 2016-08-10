package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"
	//"io/ioutil" // Needed for the json.unmarshal usage below
	"net/http"
)

func searchDeadShows(numberOfResults int, startPage int) []ArchiveDoc {

	var searchURL = fmt.Sprint("https://archive.org/advancedsearch.php?"+
		"q=collection%3AGratefulDead&fl%5B%5D=identifier&sort%5B%5D=identifier+"+
		"desc&rows=", numberOfResults, "&page=", startPage, "&output=json")

	// Print search URL to the console
	fmt.Println(searchURL)

	r, _ := http.Get(searchURL)
	defer r.Body.Close()

	var searchResponse ArchiveSearchResponse
	// Should this be a look until EOF? Or am I sure I have all the response?
	// Is that what the Body.Close call gets me?
	json.NewDecoder(r.Body).Decode(&searchResponse)

	// Print search result details to the console
	fmt.Println("Rows found: ", searchResponse.Response.NumFound)
	fmt.Println("Rows returned: ", searchResponse.ResponseHeader.Params.Rows)
	fmt.Println("Starting record: ", searchResponse.Response.Start)

	var docs = []ArchiveDoc{}

	for _, doc := range searchResponse.Response.Docs {
		fmt.Println(doc.Identifier)
		docs = append(docs, doc)
	}

	return docs
}

// TODO: NO error handling. Should probably do something about that.
func main() {

	numResults := flag.Int("numResults", 50, "Number of results to return")
	startPage := flag.Int("startPage", 0, "Results page to start with")
	flag.Parse()

	fmt.Println("numResults: ", *numResults)
	fmt.Println("startPage: ", *startPage)

	// Print show details to the console
	results := searchDeadShows(*numResults, *startPage)
	for _, doc := range results {
		fmt.Println(doc.Identifier)
		showURL := "http://archive.org/details/" + doc.Identifier + "?output=json"
		fmt.Println(showURL)
		r, _ := http.Get(showURL)
		defer r.Body.Close()

		showJSON, _ := ioutil.ReadAll(r.Body)

		var showResponse DeadShow
		// Should this be a look until EOF? Or am I sure I have all the response?
		// Is that what the Body.Close call gets me?
		//json.NewDecoder(r.Body).Decode(&showResponse)

		// Alternative solution using Unmarshal
		// show, _ := ioutil.ReadAll(s.Body)
		json.Unmarshal(showJSON, &showResponse)

		var deadFiles DeadShowFiles
		json.Unmarshal(showJSON, &deadFiles)
		for k, v := range deadFiles.Files {
			fmt.Println(k)
			fmt.Println(v)
		}

		// :( All these will panic if value is nil
		fmt.Println("Server: ", showResponse.Server)
		fmt.Println("Date: ", showResponse.Metadata.Date[0])
		t, _ := time.Parse("2006-01-02", showResponse.Metadata.Date[0])
		fmt.Println("Parsed Date: ", t.Format("2006-01-02"))
		fmt.Println("Venue: ", showResponse.Metadata.Venue[0])
		fmt.Println("Backup Location: ", showResponse.Metadata.BackupLocation[0])

		// Dump the whole showResponse to console
		//fmt.Printf("%+v\n", showResponse)

		fmt.Scanln()
	}
}
