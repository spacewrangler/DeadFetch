package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	//"io/ioutil" // Needed for the json.unmarshal usage below
	"net/http"
	"time"

	"gopkg.in/olivere/elastic.v3"
)

// TODO: NO error handling. Should probably do something about that.
func main() {

	numResults := flag.Int("numResults", 50, "Number of results to return")
	startPage := flag.Int("startPage", 0, "Results page to start with")
	flag.Parse()

	fmt.Println("numResults: ", *numResults)
	fmt.Println("startPage: ", *startPage)

	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Trace request and response details like this
	//client.SetTracer(log.New(os.Stdout, "", 0))

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	exists, err := client.IndexExists("deadshows").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create an index
		// TODO increase 1K field limit
		createIndex, err := client.CreateIndex("deadshows").Do()
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// Wait for user input to continue
	fmt.Scanln()

	var searchURL = fmt.Sprint("https://archive.org/advancedsearch.php?"+
		"q=collection%3AGratefulDead&fl%5B%5D=identifier&sort%5B%5D=identifier+"+
		"desc&rows=", *numResults, "&page=", *startPage, "&output=json")

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
	var id int64
	for _, doc := range searchResponse.Response.Docs {
		fmt.Println(doc.Identifier)
		var showURL = "http://archive.org/details/" + doc.Identifier + "?output=json"
		fmt.Println(showURL)
		r, _ := http.Get(showURL)
		defer r.Body.Close()

		showJSON, _ := ioutil.ReadAll(r.Body)
		//fmt.Println(string(showJSON))

		put1, err := client.Index().
			Index("deadshows").
			Type("deadshow").
			Id("id").
			BodyString(string(showJSON)).
			Do()
		if err != nil {
			// Handle error
			panic(err)
		}

		fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

		get1, err := client.Get().
			Index("deadshows").
			Type("deadshow").
			Id(string(id)).
			Do()
		if err != nil {
			// Handle error
			panic(err)
		}

		if get1.Found {
			fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
		}

		var showResponse DeadShow
		// Should this be a look until EOF? Or am I sure I have all the response?
		// Is that what the Body.Close call gets me?
		json.NewDecoder(r.Body).Decode(&showResponse)

		// Alternative solution using Unmarshal
		// show, _ := ioutil.ReadAll(s.Body)
		// json.Unmarshal(show, &showResponse)

		// :( All these will panic if value is nil
		fmt.Println("Server: ", showResponse.Server)
		fmt.Println("Date: ", showResponse.Metadata.Date[0])
		t, _ := time.Parse("2006-01-02", showResponse.Metadata.Date[0])
		fmt.Println("Parsed Date: ", t.Format("2006-01-02"))
		fmt.Println("Venue: ", showResponse.Metadata.Venue[0])

		id++

		fmt.Scanln()
	}
}
