package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	//"io/ioutil" // Needed for the json.unmarshal usage below
	"net/http"
)

// TODO: NO error handling. Should probably do something about that.
func main() {
	fmt.Println("***********************")
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
		json.Unmarshal(showJSON, &showResponse)
		b, _ := (json.MarshalIndent(showResponse.Details, "", "   "))
		println(string(b))
		fmt.Println("***********************")

		fmt.Scanln()
	}
}
