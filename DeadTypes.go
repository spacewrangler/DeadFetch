package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

func unmarshalDeadShowDetails(raw *DeadShowRaw, show *DeadShow) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s: %s", r, debug.Stack())
		}
	}()
	// TODO Fix format: should include time
	if raw.Metadata.Addeddate != nil {
		show.Details.AddedDate, _ = time.Parse("2006-01-02", raw.Metadata.Addeddate[0])
	}
	// TODO convert to float
	//show.Details.AverageReview = raw.Reviews.Info.AvgRating
	if raw.Metadata.BackupLocation != nil {
		show.Details.BackupLocation = raw.Metadata.BackupLocation[0]
	}
	show.Details.Collection = raw.Metadata.Collection
	show.Details.CollectionTitle = raw.Misc.CollectionTitle
	if raw.Metadata.Coverage != nil {
		show.Details.Coverage = raw.Metadata.Coverage[0]
	}
	if raw.Metadata.Creator != nil {
		show.Details.Creator = raw.Metadata.Creator[0]
	}

	show.Details.Curation = raw.Metadata.Curation
	if raw.Metadata.Date != nil {
		show.Details.Date, _ = time.Parse("2006-01-02", raw.Metadata.Date[0])
	}
	if raw.Metadata.Description != nil {
		show.Details.Description = raw.Metadata.Description[0]
	}
	show.Details.DownloadsMonth = raw.Item.Month
	show.Details.DownloadsWeek = raw.Item.Week
	show.Details.ImageLink = raw.Misc.Image
	if raw.Metadata.Lineage != nil {
		show.Details.Lineage = raw.Metadata.Lineage[0]
	}
	show.Details.MD5s = raw.Metadata.Md5S
	show.Details.Mediatype = raw.Metadata.Mediatype
	show.Details.Notes = raw.Metadata.Notes
	// TODO Fix format: should include time
	if raw.Metadata.Publicdate != nil {
		show.Details.PublicDate, _ = time.Parse("2006-01-02", raw.Metadata.Publicdate[0])
	}
	show.Details.ReviewCount = raw.Reviews.Info.NumReviews
	// TODO runtime string format is not standard - write parsing function
	//show.Details.RunTime =
	if raw.Metadata.Source != nil {
		show.Details.Source = raw.Metadata.Source[0]
	}
	if raw.Metadata.Subject != nil {
		show.Details.Subject = raw.Metadata.Subject[0]
	}
	if raw.Metadata.Taper != nil {
		show.Details.Taper = raw.Metadata.Taper[0]
	}
	if raw.Metadata.Title != nil {
		show.Details.Title = raw.Metadata.Title[0]
	}
	show.Details.TotalDownloads = raw.Item.Downloads
	if raw.Metadata.Transferer != nil {
		show.Details.Transferer = raw.Metadata.Transferer[0]
	}
	show.Details.Type = raw.Metadata.Type
	//TODO Fix format: should include time
	if raw.Metadata.Updatedate != nil {
		show.Details.UpdateDate, _ = time.Parse("2006-01-02", raw.Metadata.Updatedate[0])
	}
	if raw.Metadata.Updater != nil {
		show.Details.Updater = raw.Metadata.Updater[0]
	}
	if raw.Metadata.Venue != nil {
		show.Details.Venue = raw.Metadata.Venue[0]
	}
	if raw.Metadata.Year != nil {
		show.Details.Year, _ = strconv.ParseUint(raw.Metadata.Year[0], 0, 64)
	}

	return nil
}

func unmarshalDeadShowReviews(raw *DeadShowRaw, show *DeadShow) error {

	// TODO: convert to float
	//show.Details.AverageReview = raw.Reviews.Info.AvgRating
	for _, r := range raw.Reviews.Reviews {
		var rev = DeadShowReview{}
		rev.Body = r.Reviewbody
		// TODO convert string to date
		//rev.Date =
		rev.Reviewer = r.Reviewer
		rev.Title = r.Reviewtitle
		// TODO convert to int
		//rev.Stars = r.Stars
		show.Reviews = append(show.Reviews, rev)
	}
	return nil
}

func unmarshalDeadShowFiles(raw *DeadShowRaw, show *DeadShow) error {
	if raw.Files != nil {
		for k, v := range raw.Files {
			var file = DeadShowFile{}

			file.Album = v.Album
			// TODO convert
			//file.Bitrate = v.Bitrate
			file.CRC32 = v.Crc32
			file.Creator = v.Creator
			file.Format = v.Format
			// TODO convert
			//file.Height = v.Height
			// TODO convert
			//file.Length = v.Length
			file.MD5 = v.Md5
			// TODO convert
			//file.MTime = v.Mtime
			file.Name = k
			file.Original = v.Original
			file.SHA1 = v.Sha1
			// TODO convert
			//file.Size = v.Size
			show.Files = append(show.Files, file)
		}
	}
	return nil
}

func (ds *DeadShow) UnmarshalJSON(data []byte) error {
	var temp DeadShowRaw
	json.Unmarshal(data, &temp)

	ds.Identifier = temp.Metadata.Identifier[0]
	ds.Server = temp.Server
	ds.Directory = temp.Dir

	unmarshalDeadShowDetails(&temp, ds)
	unmarshalDeadShowReviews(&temp, ds)
	unmarshalDeadShowFiles(&temp, ds)

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
	Stars    uint8
	Body     string
}

type DeadShowFile struct {
	Name     string
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
		docs = append(docs, doc)
	}

	return docs
}
