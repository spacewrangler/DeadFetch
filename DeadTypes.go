package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func unmarshalDeadShowDetails(raw *DeadShowRaw, show *DeadShow) error {
	// TODO Fix format: should include time
	show.Details.AddedDate, _ = time.Parse("2006-01-02", raw.Metadata.Addeddate[0])
	show.Details.AverageReview, _ = strconv.ParseFloat(raw.Reviews.Info.AvgRating, 64)
	show.Details.BackupLocation = raw.Metadata.BackupLocation[0]
	show.Details.Collection = raw.Metadata.Collection
	show.Details.CollectionTitle = raw.Misc.CollectionTitle
	show.Details.Coverage = raw.Metadata.Coverage[0]
	show.Details.Creator = raw.Metadata.Creator[0]
	show.Details.Curation = raw.Metadata.Curation
	show.Details.Date, _ = time.Parse("2006-01-02", raw.Metadata.Date[0])
	show.Details.Description = raw.Metadata.Description[0]
	show.Details.DownloadsMonth = raw.Item.Month
	show.Details.DownloadsWeek = raw.Item.Week
	show.Details.ImageLink = raw.Misc.Image
	show.Details.Lineage = raw.Metadata.Lineage[0]
	show.Details.MD5s = raw.Metadata.Md5S
	show.Details.Mediatype = raw.Metadata.Mediatype
	show.Details.Notes = raw.Metadata.Notes
	// TODO Fix format: should include time
	show.Details.PublicDate, _ = time.Parse("2006-01-02", raw.Metadata.Publicdate[0])
	show.Details.ReviewCount = raw.Reviews.Info.NumReviews
	// TODO runtime string format is not standard - write parsing function
	//show.Details.RunTime =
	show.Details.Source = raw.Metadata.Source[0]
	show.Details.Subject = raw.Metadata.Subject[0]
	show.Details.Taper = raw.Metadata.Taper[0]
	show.Details.Title = raw.Metadata.Title[0]
	show.Details.TotalDownloads = raw.Item.Downloads
	show.Details.Transferer = raw.Metadata.Transferer[0]
	show.Details.Type = raw.Metadata.Type
	//TODO Fix format: should include time
	show.Details.UpdateDate, _ = time.Parse("2006-01-02", raw.Metadata.Updatedate[0])
	show.Details.Updater = raw.Metadata.Updater[0]
	show.Details.Venue = raw.Metadata.Venue[0]
	show.Details.Year, _ = strconv.ParseUint(raw.Metadata.Year[0], 0, 64)

	return nil
}

func (ds *DeadShow) UnmarshalJSON(data []byte) error {
	var temp DeadShowRaw
	json.Unmarshal(data, &temp)

	ds.Identifier = temp.Metadata.Identifier[0]
	ds.Server = temp.Server
	ds.Directory = temp.Dir

	unmarshalDeadShowDetails(&temp, ds)

	return nil
}

type DeadShow struct {
	Identifier string
	Server     string
	Directory  string
	Details    DeadShowDetails
	Reviews    []DeadShowReview
	Files      []DeadShowFile
}

type DeadShowDetails struct {
	Title           string
	Creator         string
	Mediatype       []string
	Collection      []string
	Type            []string
	Description     string
	Date            time.Time
	Year            uint64
	Subject         string
	PublicDate      time.Time
	AddedDate       time.Time
	Venue           string
	Coverage        string
	Source          string
	Lineage         string
	Taper           string
	Transferer      string
	RunTime         time.Duration
	MD5s            []string
	Notes           []string
	UpdateDate      time.Time
	Updater         string
	Curation        []string
	BackupLocation  string
	ReviewCount     uint32
	AverageReview   float64
	TotalDownloads  uint32
	DownloadsWeek   uint32
	DownloadsMonth  uint32
	ImageLink       string
	CollectionTitle string
}

type DeadShowReview struct {
	Title    string
	Reviewer string
	Date     time.Time
	Stars    int
	Body     string
}

type DeadShowFile struct {
	Source   string
	Creator  string
	Title    string
	Track    uint16
	Album    string
	Bitrate  uint16
	Format   string
	Original string
	MTime    uint64
	Size     uint64
	MD5      string
	CRC32    string
	SHA1     string
	Length   time.Duration
	Height   uint16
	Width    uint16
}

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
