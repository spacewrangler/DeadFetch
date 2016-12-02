package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"gopkg.in/olivere/elastic.v5"
)

// TODO: NO error handling. Should probably do something about that.
func main() {
	LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	numResults := flag.Int("numResults", 50, "Number of results to return")
	startPage := flag.Int("startPage", 0, "Results page to start with")
	flag.Parse()

	fmt.Println("numResults: ", *numResults)
	fmt.Println("startPage: ", *startPage)

	// Setup elastic
	errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)
	client, err := elastic.NewClient(elastic.SetErrorLog(errorlog))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Trace request and response details like this
	//client.SetTracer(log.New(os.Stdout, "", 0))
	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(context.TODO())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	exists, err := client.IndexExists("deadshows").Do(context.TODO())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create an index
		// TODO increase 1K field limit
		createIndex, err := client.CreateIndex("deadshows").Do(context.TODO())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// Print show details to the console
	results := SearchDeadShows(*numResults, *startPage)
	Trace.Println("Number of results: {0}", len(results))
	for _, doc := range results {
		showURL := "http://archive.org/details/" + doc.Identifier + "?output=json"
		Trace.Println(showURL)
		r, _ := http.Get(showURL)
		defer r.Body.Close()

		showJSON, _ := ioutil.ReadAll(r.Body)

		var showResponse DeadShow
		json.Unmarshal(showJSON, &showResponse)

		show, _ := json.Marshal(showResponse.Details)
		fmt.Println(*showResponse.Details.Coverage)
		fmt.Println(convertCityToLatLng(*showResponse.Details.Coverage))

		Trace.Println("Indexing: {0}", *showResponse.Identifier)
		put1, err := client.Index().
			Index("deadshows").
			Type("deadshow").
			Id(string(*showResponse.Identifier)).
			// TODO: this has been null before. Check that show is actually retrieved.
			// If not, re-fetch
			BodyString(string(show)).
			Do(context.TODO())
		if err != nil {
			// Handle error
			panic(err)
		}

		Trace.Println("Created:", put1.Created, "Version:", put1.Version)

		// get1, err := client.Get().
		// 	Index("deadshows").
		// 	Type("deadshow").
		// 	Id(string(*showResponse.Identifier)).
		// 	Do()
		// if err != nil {
		// 	// Handle error
		// 	panic(err)
		// }

		// if get1.Found {
		// 	fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
		// }

		//fmt.Scanln()
	}
}
