package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Params contains the query parameters passed to archive.org
type archiveSearchQueryParams struct {
	Query   string `json:"q"`
	QueryIn string `json:"qin"`
	Rows    int    `json:",string"`
	Start   int
}

type archiveSearchResponseHeader struct {
	Status int
	QTime  int
	Params archiveSearchQueryParams
}

type ArchiveDoc struct {
	Identifier    string
	OaiUpdateDate []time.Time `json:"oai_updatedate"`
}

type archiveSearchResults struct {
	NumFound int
	Start    int
	Docs     []ArchiveDoc
}

type archiveSearchResponse struct {
	ResponseHeader archiveSearchResponseHeader
	Response       archiveSearchResults
}

func SearchDeadShows(numberOfResults int, startPage int) []ArchiveDoc {

	var searchURL = fmt.Sprint("https://archive.org/advancedsearch.php?"+
		"q=collection%3AGratefulDead&fl%5B%5D=identifier&sort%5B%5D=identifier+"+
		"desc&rows=", numberOfResults, "&page=", startPage, "&output=json")

	// Print search URL to the console
	fmt.Println(searchURL)

	r, _ := http.Get(searchURL)
	defer r.Body.Close()

	var searchResponse archiveSearchResponse
	// Should this be a look until EOF? Or am I sure I have all the response?
	// Is that what the Body.Close call gets me?
	json.NewDecoder(r.Body).Decode(&searchResponse)

	// Print search result details to the console
	fmt.Println("Rows found: ", searchResponse.Response.NumFound)
	fmt.Println("Rows returned: ", searchResponse.ResponseHeader.Params.Rows)
	fmt.Println("Starting record: ", searchResponse.Response.Start)

	var docs = []ArchiveDoc{}

	for _, doc := range searchResponse.Response.Docs {
		docs = append(docs, doc)
	}

	return docs
}
