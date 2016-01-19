package main

import (
    //"time"
)

// Params contains the query parameters passed to archive.org
type Params struct {
	Query   string `json:"q"`
	QueryIn string `json:"qin"`
	Rows    int    `json:",string"`
	Start   int
}

type ResponseHeader struct {
	Status int
	QTime  int
	Params Params
}

type Doc struct {
	Identifier string
}

type Response struct {
	NumFound int
	Start    int
	Docs     []Doc
}

type SearchResponse struct {
	ResponseHeader ResponseHeader
	Response       Response
}

type DeadShow struct {
    Server      string
    Metadata    DeadShowMetadata `json:"metadata"`
}

type DeadShowMetadata struct {
    Identifier  []string
    Date       []string
    UpdateDate  []string
    Venue      []string    
}
