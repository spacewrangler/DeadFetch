package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

type deadShowRaw struct {
	Server   string
	Metadata deadShowMetadataRaw `json:"metadata"`
}

type DeadShow struct {
	Server         string
	Identifier     []string
	Title          []string
	Creator        []string
	Mediatype      []string
	Collection     []string
	Type           []string
	Description    []string
	Date           []string
	Year           []string
	Subject        []string
	PublicDate     []string
	AddedDate      []string
	Venue          []string
	Coverage       []string
	Source         []string
	Lineage        []string
	Taper          []string
	Transferer     []string
	RunTime        []string
	md5s           []string
	Notes          []string
	UpdateDate     []string
	Updater        []string
	Curation       []string
	BackupLocation []string
}

func (ds *DeadShow) UnmarshalJSON(data []byte) error {
	var dsr deadShowRaw
	json.Unmarshal(data, &dsr)

	ds.Server = dsr.Server
	ds.Identifier = dsr.Metadata.Identifier
	ds.Title = dsr.Metadata.Title
	ds.Creator = dsr.Metadata.Creator
	ds.Mediatype = dsr.Metadata.Mediatype
	ds.Collection = dsr.Metadata.Collection
	ds.Type = dsr.Metadata.Type
	ds.Description = dsr.Metadata.Description
	ds.Date = dsr.Metadata.Date
	ds.Year = dsr.Metadata.Year
	ds.Subject = dsr.Metadata.Year
	ds.PublicDate = dsr.Metadata.PublicDate
	ds.AddedDate = dsr.Metadata.AddedDate
	ds.Venue = dsr.Metadata.Venue
	ds.Coverage = dsr.Metadata.Coverage
	ds.Source = dsr.Metadata.Source
	ds.Lineage = dsr.Metadata.Lineage
	ds.Taper = dsr.Metadata.Taper
	ds.Transferer = dsr.Metadata.Transferer
	ds.RunTime = dsr.Metadata.RunTime
	ds.md5s = dsr.Metadata.md5s
	ds.Notes = dsr.Metadata.Notes
	ds.UpdateDate = dsr.Metadata.UpdateDate
	ds.Updater = dsr.Metadata.Updater
	ds.Curation = dsr.Metadata.Curation
	ds.BackupLocation = dsr.Metadata.BackupLocation

	return nil
}

type deadShowMetadataRaw struct {
	Identifier     []string
	Title          []string
	Creator        []string
	Mediatype      []string
	Collection     []string
	Type           []string
	Description    []string
	Date           []string
	Year           []string
	Subject        []string
	PublicDate     []string
	AddedDate      []string
	Venue          []string
	Coverage       []string
	Source         []string
	Lineage        []string
	Taper          []string
	Transferer     []string
	RunTime        []string
	md5s           []string
	Notes          []string
	UpdateDate     []string
	Updater        []string
	Curation       []string
	BackupLocation []string `json:"backup_location"`
}

type deadShowFilesRaw struct {
	Files map[string]deadShowFileRaw
}

type deadShowFileRaw struct {
	Name   string
	Source string
	Format string
}

type DeadShowFile struct {
	// TODO add show ID
	//ShowIdentifier string
	Name   string
	Source string
	Format string
}

type DeadShowFiles struct {
	Files []DeadShowFile
}

func (dsf *DeadShowFiles) UnmarshalJSON(data []byte) error {
	// These var names are confusing: too many "file"
	var deadFiles deadShowFilesRaw
	var deadFile DeadShowFile
	json.Unmarshal(data, &deadFiles)
	for k, v := range deadFiles.Files {
		deadFile.Name = strings.TrimPrefix(k, "/")
		deadFile.Source = v.Source
		deadFile.Format = v.Format
		dsf.Files = append(dsf.Files, deadFile)
	}

	return nil
}

func searchDeadShows(numberOfResults int, startPage int) []ArchiveDoc {

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
		fmt.Println(doc.Identifier)
		docs = append(docs, doc)
	}

	return docs
}
