package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type GeoLocation struct {
	AddressComponents []AddressComponents `json:"address_components,omitempty"`
	FormattedAddress  string              `json:"formatted_address,omitempty"`
	Geometry          Geometry            `json:"geometry,omitempty"`
	PlaceID           string              `json:"place_id,omitempty"`
	Types             []string            `json:"types,omitempty"`
}
type AddressComponents struct {
	LongName  string   `json:"long_name,omitempty"`
	ShortName string   `json:"short_name,omitempty"`
	Types     []string `json:"types,omitempty"`
}
type Northeast struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}
type Southwest struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}
type Bounds struct {
	Northeast Northeast `json:"northeast,omitempty"`
	Southwest Southwest `json:"southwest,omitempty"`
}
type Location struct {
	Lat float64 `json:"lat,omitempty"`
	Lng float64 `json:"lng,omitempty"`
}
type Viewport struct {
	Northeast Northeast `json:"northeast,omitempty"`
	Southwest Southwest `json:"southwest,omitempty"`
}
type Geometry struct {
	Bounds       Bounds   `json:"bounds,omitempty"`
	Location     Location `json:"location,omitempty"`
	LocationType string   `json:"location_type,omitempty"`
	Viewport     Viewport `json:"viewport,omitempty"`
}

type DeadShow struct {
	Identifier *string          `json:",omitempty"`
	Server     *string          `json:",omitempty"`
	Directory  *string          `json:",omitempty"`
	Details    DeadShowDetails  `json:",omitempty"`
	Reviews    []DeadShowReview `json:",omitempty"`
	Files      []DeadShowFile   `json:",omitempty"`
}

type DeadShowDetails struct {
	Title           *string        `json:",omitempty"`
	Creator         *string        `json:",omitempty"`
	Mediatype       []string       `json:",omitempty"`
	Collection      []string       `json:",omitempty"`
	Type            []string       `json:",omitempty"`
	Description     *string        `json:",omitempty"`
	Date            *time.Time     `json:",omitempty"`
	Year            *uint64        `json:",omitempty"`
	Subject         *string        `json:",omitempty"`
	PublicDate      *time.Time     `json:",omitempty"`
	AddedDate       *time.Time     `json:",omitempty"`
	Venue           *string        `json:",omitempty"`
	Location        *string        `json:",omitempty"`
	GeoLocation     *string        `json:",omitempty"`
	Source          *string        `json:",omitempty"`
	Lineage         *string        `json:",omitempty"`
	Taper           *string        `json:",omitempty"`
	Transferer      *string        `json:",omitempty"`
	RunTime         *time.Duration `json:",omitempty"`
	Notes           []string       `json:",omitempty"`
	UpdateDate      *time.Time     `json:",omitempty"`
	Updater         *string        `json:",omitempty"`
	Curation        []string       `json:",omitempty"`
	BackupLocation  *string        `json:",omitempty"`
	ReviewCount     *uint32        `json:",omitempty"`
	AverageReview   *float64       `json:",omitempty"`
	TotalDownloads  *uint32        `json:",omitempty"`
	DownloadsWeek   *uint32        `json:",omitempty"`
	DownloadsMonth  *uint32        `json:",omitempty"`
	ImageLink       *string        `json:",omitempty"`
	CollectionTitle *string        `json:",omitempty"`
	SetList         []string       `json:",omitempty"`
}

type DeadShowReview struct {
	Title    *string    `json:",omitempty"`
	Reviewer *string    `json:",omitempty"`
	Date     *time.Time `json:",omitempty"`
	Stars    *uint64    `json:",omitempty"`
	Body     *string    `json:",omitempty"`
}

type DeadShowFile struct {
	Name     *string        `json:",omitempty"`
	Source   *string        `json:",omitempty"`
	Creator  *string        `json:",omitempty"`
	Title    *string        `json:",omitempty"`
	Track    *uint64        `json:",omitempty"`
	Album    *string        `json:",omitempty"`
	Bitrate  *uint64        `json:",omitempty"`
	Format   *string        `json:",omitempty"`
	Original *string        `json:",omitempty"`
	MTime    *uint64        `json:",omitempty"`
	Size     *uint64        `json:",omitempty"`
	MD5      *string        `json:",omitempty"`
	CRC32    *string        `json:",omitempty"`
	SHA1     *string        `json:",omitempty"`
	Length   *time.Duration `json:",omitempty"`
	Height   *uint64        `json:",omitempty"`
	Width    *uint64        `json:",omitempty"`
}

func geocodingResultToGeoLocation(gc maps.GeocodingResult) GeoLocation {
	var gl GeoLocation

	for i := range gc.AddressComponents {
		var ac AddressComponents
		ac.LongName = gc.AddressComponents[i].LongName
		ac.ShortName = gc.AddressComponents[i].ShortName
		ac.Types = gc.AddressComponents[i].Types
		gl.AddressComponents = append(gl.AddressComponents, ac)
	}
	gl.FormattedAddress = gc.FormattedAddress
	gl.Geometry.Bounds.Northeast.Lat = gc.Geometry.Bounds.NorthEast.Lat
	gl.Geometry.Bounds.Northeast.Lng = gc.Geometry.Bounds.NorthEast.Lng
	gl.Geometry.Bounds.Southwest.Lat = gc.Geometry.Bounds.NorthEast.Lat
	gl.Geometry.Bounds.Southwest.Lng = gc.Geometry.Bounds.NorthEast.Lng
	gl.Geometry.Location.Lat = gc.Geometry.Location.Lat
	gl.Geometry.Location.Lng = gc.Geometry.Location.Lng
	gl.Geometry.LocationType = gc.Geometry.LocationType
	gl.Geometry.Viewport.Northeast.Lat = gc.Geometry.Viewport.NorthEast.Lat
	gl.Geometry.Viewport.Northeast.Lng = gc.Geometry.Viewport.NorthEast.Lng
	gl.Geometry.Viewport.Southwest.Lat = gc.Geometry.Viewport.SouthWest.Lat
	gl.Geometry.Viewport.Southwest.Lng = gc.Geometry.Viewport.SouthWest.Lng
	gl.PlaceID = gc.PlaceID
	gl.Types = gc.Types

	return gl
}

func cityToGeoLocation(city string) GeoLocation {
	var gl GeoLocation

	q := elastic.NewMatchQuery("formatted_address", city).
		Operator("AND").
		Fuzziness("AUTO")

	sr, err := client.Search().
		Index("geolocation").
		Query(q).
		From(0).
		Size(1).
		Do(context.TODO())

	// Not found in elasticsearch
	if sr.TotalHits() == 0 {

		// Fetch from Google
		// TODO: hoist this out so I'm not reading the file and creating connection each time
		dat, err := ioutil.ReadFile("api.key")
		s := string(dat)
		if err != nil {
			panic(err)
		}

		c, err := maps.NewClient(maps.WithAPIKey(s))
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}
		r := &maps.GeocodingRequest{
			Address: city,
		}
		resp, err := c.Geocode(context.Background(), r)
		if err != nil {
			panic(err)
		}
		if len(resp) != 0 {
			gl = geocodingResultToGeoLocation(resp[0])
			put1, err := client.Index().
				Index("geolocation").
				Id(gl.PlaceID).
				BodyJson(gl).
				Do(context.Background())
			if err != nil {
				// Handle error
				panic(err)
			}
			log.Printf("Index result: %s, status: %d\n", put1.Result, put1.Status)
		}
	} else {
		json.Unmarshal(sr.Hits.Hits[0].Source, &gl)
		if err != nil {
			panic(err)
		}
	}
	return gl
}

func unmarshalDeadShowDetails(raw *DeadShowRaw, show *DeadShow) error {

	if raw.Metadata.Addeddate != nil {
		if raw.Metadata.Addeddate[0] == "" {
			show.Details.AddedDate = nil
		} else {
			t, _ := time.Parse("2006-01-02 15:04:05", raw.Metadata.Addeddate[0])
			show.Details.AddedDate = &t
		}
	}

	if raw.Reviews.Info.AvgRating != nil {
		if *raw.Reviews.Info.AvgRating == "" {
			show.Details.AverageReview = nil
		} else {
			a, _ := strconv.ParseFloat(*raw.Reviews.Info.AvgRating, 64)
			show.Details.AverageReview = &a
		}
	}

	if raw.Metadata.BackupLocation != nil {
		if raw.Metadata.BackupLocation[0] == "" {
			show.Details.BackupLocation = nil
		} else {
			show.Details.BackupLocation = &raw.Metadata.BackupLocation[0]
		}
	}

	show.Details.Collection = raw.Metadata.Collection

	if raw.Misc.CollectionTitle != nil {
		if *raw.Misc.CollectionTitle == "" {
			show.Details.Collection = nil
		} else {
			show.Details.CollectionTitle = raw.Misc.CollectionTitle
		}
	}

	if raw.Metadata.Coverage != nil {
		if raw.Metadata.Coverage[0] == "" {
			show.Details.Location = nil
		} else {
			if config.GeoLookup == true {
				gl := cityToGeoLocation(raw.Metadata.Coverage[0])
				show.Details.Location = &gl.FormattedAddress
				s := fmt.Sprintf("%f,%f", gl.Geometry.Location.Lat, gl.Geometry.Location.Lng)
				show.Details.GeoLocation = &s
			} else {
				show.Details.GeoLocation = nil
			}
		}
	}

	if raw.Metadata.Creator != nil {
		if raw.Metadata.Creator[0] == "" {
			show.Details.Creator = nil
		} else {
			show.Details.Creator = &raw.Metadata.Creator[0]
		}
	}

	show.Details.Curation = raw.Metadata.Curation

	if raw.Metadata.Date != nil {
		if raw.Metadata.Date[0] == "" {
			show.Details.Date = nil
		} else {
			t, _ := time.Parse("2006-01-02", raw.Metadata.Date[0])
			show.Details.Date = &t
		}
	}

	if raw.Metadata.Description != nil {
		if raw.Metadata.Description[0] == "" {
			show.Details.Description = nil
		} else {
			show.Details.Description = &raw.Metadata.Description[0]
		}
	}

	show.Details.DownloadsMonth = raw.Item.Month
	show.Details.DownloadsWeek = raw.Item.Week

	if raw.Misc.Image != nil {
		if *raw.Misc.Image == "" {
			show.Details.ImageLink = nil
		} else {
			show.Details.ImageLink = raw.Misc.Image
		}
	}

	if raw.Metadata.Lineage != nil {
		if raw.Metadata.Lineage[0] == "" {
			show.Details.Lineage = nil
		} else {
			show.Details.Lineage = &raw.Metadata.Lineage[0]
		}
	}

	show.Details.Mediatype = raw.Metadata.Mediatype
	show.Details.Notes = raw.Metadata.Notes

	if raw.Metadata.Publicdate != nil {
		t, _ := time.Parse("2006-01-02 15:04:05", raw.Metadata.Publicdate[0])
		show.Details.PublicDate = &t
	}

	show.Details.ReviewCount = raw.Reviews.Info.NumReviews
	// TODO: runtime string format is not standard - write parsing function
	//show.Details.RunTime =

	if raw.Metadata.Source != nil {
		if raw.Metadata.Source[0] == "" {
			show.Details.Source = nil
		} else {
			show.Details.Source = &raw.Metadata.Source[0]
		}
	}

	if raw.Metadata.Subject != nil {
		if raw.Metadata.Subject[0] == "" {
			show.Details.Subject = nil
		} else {
			show.Details.Subject = &raw.Metadata.Subject[0]
		}
	}

	if raw.Metadata.Taper != nil {
		if raw.Metadata.Taper[0] == "" {
			show.Details.Taper = nil
		} else {
			show.Details.Taper = &raw.Metadata.Taper[0]
		}
	}
	if raw.Metadata.Title != nil {
		if raw.Metadata.Title[0] == "" {
			show.Details.Title = nil
		} else {
			show.Details.Title = &raw.Metadata.Title[0]
		}
	}

	show.Details.TotalDownloads = raw.Item.Downloads

	if raw.Metadata.Transferer != nil {
		if raw.Metadata.Transferer[0] == "" {
			show.Details.Transferer = nil
		} else {
			show.Details.Transferer = &raw.Metadata.Transferer[0]
		}
	}
	show.Details.Type = raw.Metadata.Type
	// TODO: This is a list of dates - should get them all
	if raw.Metadata.Updatedate != nil {
		if raw.Metadata.Updatedate[0] == "" {
			show.Details.UpdateDate = nil
		} else {
			t, _ := time.Parse("2006-01-02 15:04:05", raw.Metadata.Updatedate[0])
			show.Details.UpdateDate = &t
		}
	}

	if raw.Metadata.Updater != nil {
		if raw.Metadata.Updater[0] == "" {
			show.Details.UpdateDate = nil
		} else {
			show.Details.Updater = &raw.Metadata.Updater[0]
		}
	}

	if raw.Metadata.Venue != nil {
		if raw.Metadata.Venue[0] == "" {
			show.Details.Venue = nil
		} else {
			show.Details.Venue = &raw.Metadata.Venue[0]
		}
	}

	if raw.Metadata.Year != nil {
		if raw.Metadata.Year[0] == "" {
			show.Details.Year = nil
		} else {
			u, _ := strconv.ParseUint(raw.Metadata.Year[0], 10, 64)
			show.Details.Year = &u
		}
	}

	return nil
}

func unmarshalDeadShowReviews(raw *DeadShowRaw, show *DeadShow) error {

	for _, r := range raw.Reviews.Reviews {
		var rev = DeadShowReview{}

		if r.Reviewbody != nil {
			if *r.Reviewbody == "" {
				rev.Body = nil
			} else {
				rev.Body = r.Reviewbody
			}
		}

		if r.Reviewdate != nil {
			if *r.Reviewdate == "" {
				rev.Date = nil
			} else {
				t, _ := time.Parse("2006-01-02 15:04:05", *r.Reviewdate)
				rev.Date = &t
			}
		}

		if r.Reviewer != nil {
			if *r.Reviewer == "" {
				rev.Reviewer = nil
			} else {
				rev.Reviewer = r.Reviewer
			}
		}

		if r.Reviewtitle != nil {
			if *r.Reviewtitle == "" {
				rev.Title = nil
			} else {
				rev.Title = r.Reviewtitle
			}
		}

		if r.Stars != nil {
			if *r.Stars == "" {
				rev.Stars = nil
			} else {
				u, _ := strconv.ParseUint(*r.Stars, 10, 64)
				rev.Stars = &u
			}
		}

		show.Reviews = append(show.Reviews, rev)
	}
	return nil
}

func unmarshalDeadShowFiles(raw *DeadShowRaw, show *DeadShow) error {
	if raw.Files != nil {

		setList := make(map[int]string)

		for k, v := range raw.Files {
			var file = DeadShowFile{}

			if v.Album != nil {
				if *v.Album == "" {
					file.Album = nil
				} else {
					file.Album = v.Album
				}
			}

			if v.Bitrate != nil {
				if *v.Bitrate == "" {
					file.Bitrate = nil
				} else {
					u, _ := strconv.ParseUint(*v.Bitrate, 10, 64)
					file.Bitrate = &u
				}
			}

			if v.Crc32 != nil {
				if *v.Crc32 == "" {
					file.CRC32 = nil
				} else {
					file.CRC32 = v.Crc32
				}
			}

			if v.Creator != nil {
				if *v.Creator == "" {
					file.Creator = nil
				} else {
					file.Creator = v.Creator
				}
			}

			if v.Format != nil {
				if *v.Format == "" {
					file.Format = nil
				} else {
					file.Format = v.Format
				}
			}

			if v.Height != nil {
				if *v.Height == "" {
					file.Height = nil
				} else {
					u, _ := strconv.ParseUint(*v.Height, 10, 64)
					file.Height = &u
				}
			}

			// TODO: Parse duration
			//file.Length = v.Length

			if v.Md5 != nil {
				if *v.Md5 == "" {
					file.MD5 = nil
				} else {
					file.MD5 = v.Md5
				}
			}

			if v.Mtime != nil {
				if *v.Mtime == "" {
					file.MTime = nil
				} else {
					u, _ := strconv.ParseUint(*v.Mtime, 10, 64)
					file.MTime = &u
				}
			}

			file.Name = &k

			if v.Original != nil {
				if *v.Original == "" {
					file.Original = nil
				} else {
					file.Original = v.Original
				}
			}

			if v.Sha1 != nil {
				if *v.Sha1 == "" {
					file.SHA1 = nil
				} else {
					file.SHA1 = v.Sha1
				}
			}

			if v.Size != nil {
				if *v.Size == "" {
					file.Size = nil
				} else {
					u, _ := strconv.ParseUint(*v.Size, 10, 64)
					file.Size = &u
				}
			}

			if v.Title != nil {
				if *v.Title == "" {
					file.Title = nil
				} else {
					file.Title = v.Title
				}
			}

			if v.Track != nil {
				if *v.Track == "" {
					file.Track = nil
				} else {
					u, _ := strconv.ParseUint(*v.Track, 10, 64)
					file.Track = &u
				}
			}

			show.Files = append(show.Files, file)

			// Create the SetList
			if strings.HasSuffix(strings.ToLower(*file.Name), "mp3") {
				if file.Track != nil && file.Title != nil {
					setList[int(*file.Track)] = *file.Title
				}
			}
		}

		for i := 0; i < len(setList); i++ {
			show.Details.SetList = append(show.Details.SetList, setList[i+1])
		}
	}

	return nil
}

func (ds *DeadShow) UnmarshalJSON(data []byte) error {
	var temp DeadShowRaw
	json.Unmarshal(data, &temp)

	if temp.Metadata.Identifier != nil {
		if temp.Metadata.Identifier[0] == "" {
			ds.Identifier = nil
		} else {
			ds.Identifier = &temp.Metadata.Identifier[0]
		}
	}

	if *temp.Server == "" {
		ds.Server = nil
	} else {
		ds.Server = temp.Server
	}

	if *temp.Dir == "" {
		ds.Directory = nil
	} else {
		ds.Directory = temp.Dir
	}

	unmarshalDeadShowDetails(&temp, ds)
	unmarshalDeadShowReviews(&temp, ds)
	unmarshalDeadShowFiles(&temp, ds)

	return nil
}
