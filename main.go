package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	//"io/ioutil" // Needed for the json.unmarshal usage below
	"net/http"
)

// TODO: NO error handling. Should probably do something about that.
func main() {
	LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

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
		Trace.Println(showURL)
		r, _ := http.Get(showURL)
		defer r.Body.Close()

		showJSON, _ := ioutil.ReadAll(r.Body)

		var showResponse DeadShow
		// Should this be a look until EOF? Or am I sure I have all the response?
		// Is that what the Body.Close call gets me?
		//json.NewDecoder(r.Body).Decode(&showResponse)

		json.Unmarshal(showJSON, &showResponse)

		// :( All these will panic if value is nil
		fmt.Println("Server: ", showResponse.Server)
		fmt.Println("Date: ", showResponse.Date[0])
		t, _ := time.Parse("2006-01-02", showResponse.Date[0])
		fmt.Println("Parsed Date: ", t.Format("2006-01-02"))
		fmt.Println("Venue: ", showResponse.Venue[0])
		fmt.Println("Backup Location: ", showResponse.BackupLocation[0])

		// Dump the whole showResponse to console
		//fmt.Printf("%+v\n", showResponse)

		var showFiles DeadShowFiles
		json.Unmarshal(showJSON, &showFiles)
		fmt.Println(showFiles.Files)

		fmt.Scanln()
	}
}
