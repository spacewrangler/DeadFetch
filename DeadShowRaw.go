package main

type DeadShowRaw struct {
	Server   string `json:"server"`
	Dir      string `json:"dir"`
	Metadata struct {
		Identifier     []string `json:"identifier"`
		Title          []string `json:"title"`
		Creator        []string `json:"creator"`
		Mediatype      []string `json:"mediatype"`
		Collection     []string `json:"collection"`
		Type           []string `json:"type"`
		Description    []string `json:"description"`
		Date           []string `json:"date"`
		Year           []string `json:"year"`
		Subject        []string `json:"subject"`
		Publicdate     []string `json:"publicdate"`
		Addeddate      []string `json:"addeddate"`
		Venue          []string `json:"venue"`
		Coverage       []string `json:"coverage"`
		Source         []string `json:"source"`
		Lineage        []string `json:"lineage"`
		Taper          []string `json:"taper,omitempty"`
		Transferer     []string `json:"transferer"`
		Runtime        []string `json:"runtime"`
		Md5S           []string `json:"md5s,-"`
		Notes          []string `json:"notes"`
		Updatedate     []string `json:"updatedate"`
		Updater        []string `json:"updater"`
		Curation       []string `json:"curation"`
		BackupLocation []string `json:"backup_location"`
	} `json:"metadata"`
	Reviews struct {
		Info struct {
			NumReviews uint32 `json:"num_reviews"`
			AvgRating  string `json:"avg_rating"`
		} `json:"info"`
		Reviews []struct {
			Reviewbody  string `json:"reviewbody"`
			Reviewtitle string `json:"reviewtitle"`
			Reviewer    string `json:"reviewer"`
			Reviewdate  string `json:"reviewdate"`
			Stars       string `json:"stars"`
		} `json:"reviews"`
	} `json:"reviews"`
	Files map[string]struct {
		Source   string `json:"source"`
		Creator  string `json:"creator"`
		Title    string `json:"title"`
		Track    string `json:"track"`
		Album    string `json:"album"`
		Bitrate  string `json:"bitrate"`
		Length   string `json:"length"`
		Format   string `json:"format"`
		Original string `json:"original"`
		Mtime    string `json:"mtime"`
		Size     string `json:"size"`
		Md5      string `json:"md5"`
		Crc32    string `json:"crc32"`
		Sha1     string `json:"sha1"`
		Height   string `json:"height"`
		Width    string `json:"width"`
	} `json:"files"`
	Misc struct {
		Image           string `json:"image"`
		CollectionTitle string `json:"collection-title"`
	} `json:"misc"`
	Item struct {
		Downloads uint32 `json:"downloads"`
		Week      uint32 `json:"week"`
		Month     uint32 `json:"month"`
	} `json:"item"`
}
