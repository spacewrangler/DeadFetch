package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
)

type Config struct {
	GeoLookup            bool
	IndexToElasticSearch bool
	ElasticsearchURL     string
	NumberOfResults      int
	StartPage            int
}

var config *Config
var client *elastic.Client

func init() {
	const (
		elasticsearchURLDefault = ""
		elasticsearchURLDesc    = "Elasticsearch URL to index Dead shows"
		numberOfResultsDefault  = 50
		numberOfResultsDesc     = "Number of results to return from Archive search"
		startPageDefault        = 0
		startPageDesc           = "Search result page to start with"
		geoLookupDefault        = false
		geoLookupDesc           = "Lookup lat/long coordinates for shows"
	)

	config = &Config{}

	flag.IntVar(&config.NumberOfResults, "numresults", numberOfResultsDefault, numberOfResultsDesc)
	flag.IntVar(&config.StartPage, "startpage", startPageDefault, startPageDesc)
	flag.StringVar(&config.ElasticsearchURL, "es-url", elasticsearchURLDefault, elasticsearchURLDesc)
	flag.BoolVar(&config.GeoLookup, "geolookup", geoLookupDefault, geoLookupDesc)
}

func initElasticsearch(url *string) {
	// Setup elastic
	//errorlog := log.New(os.Stdout, "APP ", log.LstdFlags)
	file, err := os.OpenFile("elastic.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}

	client, err = elastic.NewClient(
		elastic.SetErrorLog(log.New(file, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(file, "ELASTIC ", log.LstdFlags)),
		elastic.SetTraceLog(log.New(file, "ELASTIC ", log.LstdFlags)),
		elastic.SetURL(*url))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Trace request and response details like this
	//client.SetTracer(log.New(os.Stdout, "", 0))
	// Ping the Elasticsearch server to get e.g. the version number
	_, _, err = client.Ping(config.ElasticsearchURL).Do(context.TODO())
	if err != nil {
		// Handle error
		panic(err)
	}

	//fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	exists, err := client.IndexExists("deadshows").Do(context.TODO())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create an index
		// TODO: increase 1K field limit
		createIndex, err := client.CreateIndex("deadshows").Do(context.TODO())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
}

func indexDeadShow(show *DeadShow) {
	Trace.Println("Indexing:", *show.Identifier)

	details, _ := json.Marshal(show.Details)
	put1, err := client.Index().
		Index("deadshows").
		// Type("deadshow").
		Id(string(*show.Identifier)).
		// TODO: this has been null before. Check that show is actually retrieved.
		// If not, re-fetch
		BodyString(string(details)).
		Do(context.TODO())
	if err != nil {
		// Handle error
		panic(err)
	}

	Trace.Println("Status:", put1.Status, "Version:", put1.Version)
}

func fetchDeadShow(id string) DeadShow {
	showURL := "http://archive.org/details/" + id + "?output=json"
	Trace.Println(showURL)
	r, _ := http.Get(showURL)
	defer r.Body.Close()

	showJSON, _ := ioutil.ReadAll(r.Body)

	var showResponse DeadShow
	json.Unmarshal(showJSON, &showResponse)

	if showResponse.Details.Location != nil {
		fmt.Println(*showResponse.Details.Location)
	} else {
		fmt.Println("No location given for show")
	}
	if showResponse.Details.GeoLocation != nil {
		fmt.Println(*showResponse.Details.GeoLocation)
	} else {
		fmt.Println("No Lat/Long available for show")
	}

	return showResponse
}

func processResults(docs []ArchiveDoc, wg *sync.WaitGroup) {
	var show DeadShow

	defer wg.Done()

	for _, doc := range docs {
		show = fetchDeadShow(doc.Identifier)

		if config.IndexToElasticSearch {
			indexDeadShow(&show)

		}
	}
}

func main() {
	LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)

	flag.Parse()
	if config.ElasticsearchURL == "" {
		config.IndexToElasticSearch = false
	} else {
		config.IndexToElasticSearch = true
	}

	fmt.Println("numResults:", config.NumberOfResults)
	fmt.Println("startPage:", config.StartPage)
	fmt.Println("elasticsearch URL:", config.ElasticsearchURL)
	fmt.Println("geolookup:", config.GeoLookup)
	fmt.Println("index to ES:", config.IndexToElasticSearch)

	initElasticsearch(&config.ElasticsearchURL)

	// Print show details to the console
	results := SearchDeadShows(config.NumberOfResults, config.StartPage)
	Trace.Println("Number of results:", len(results))

	var divided [][]ArchiveDoc
	numCPU := 8

	chunkSize := (len(results) + numCPU - 1) / numCPU

	for i := 0; i < len(results); i += chunkSize {
		end := i + chunkSize

		if end > len(results) {
			end = len(results)
		}

		divided = append(divided, results[i:end])
	}
	var wg sync.WaitGroup
	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go processResults(divided[i], &wg)
	}
	wg.Wait()
}
