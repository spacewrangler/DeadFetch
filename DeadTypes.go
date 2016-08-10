package main

import (
	"time"
)

// Params contains the query parameters passed to archive.org
type ArchiveSearchQueryParams struct {
	Query   string `json:"q"`
	QueryIn string `json:"qin"`
	Rows    int    `json:",string"`
	Start   int
}

type ArchiveSearchResponseHeader struct {
	Status int
	QTime  int
	Params ArchiveSearchQueryParams
}

type ArchiveDoc struct {
	Identifier    string
	OaiUpdateDate []time.Time `json:"oai_updatedate"`
}

type ArchiveSearchResults struct {
	NumFound int
	Start    int
	Docs     []ArchiveDoc
}

type ArchiveSearchResponse struct {
	ResponseHeader ArchiveSearchResponseHeader
	Response       ArchiveSearchResults
}

type DeadShow struct {
	Server   string
	Metadata DeadShowMetadata `json:"metadata"`
}

type DeadShowMetadata struct {
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

type DeadShowFiles struct {
	Files map[string]DeadShowFile
}

type DeadShowFile struct {
	Source string
	Format string
}
