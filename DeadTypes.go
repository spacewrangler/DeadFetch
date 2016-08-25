package main

import (
	"encoding/json"
	"strconv"
	"time"
)

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

	if *raw.Misc.CollectionTitle == "" {
		show.Details.Collection = nil
	} else {
		show.Details.CollectionTitle = raw.Misc.CollectionTitle
	}

	if raw.Metadata.Coverage != nil {
		if raw.Metadata.Coverage[0] == "" {
			show.Details.Coverage = nil
		} else {
			show.Details.Coverage = &raw.Metadata.Coverage[0]
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

	if *raw.Misc.Image == "" {
		show.Details.ImageLink = nil
	} else {
		show.Details.ImageLink = raw.Misc.Image
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
	// TODO runtime string format is not standard - write parsing function
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
	// TODO This is a list of dates - should get them all
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
			u, _ := strconv.ParseUint(raw.Metadata.Year[0], 0, 64)
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
				u, _ := strconv.ParseUint(*r.Stars, 0, 64)
				rev.Stars = &u
			}
		}

		show.Reviews = append(show.Reviews, rev)
	}
	return nil
}

func unmarshalDeadShowFiles(raw *DeadShowRaw, show *DeadShow) error {
	if raw.Files != nil {
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
					u, _ := strconv.ParseUint(*v.Bitrate, 0, 64)
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
					u, _ := strconv.ParseUint(*v.Height, 0, 64)
					file.Height = &u
				}
			}

			// TODO Parse duration
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
					u, _ := strconv.ParseUint(*v.Mtime, 0, 64)
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
					u, _ := strconv.ParseUint(*v.Size, 0, 64)
					file.Size = &u
				}
			}

			show.Files = append(show.Files, file)
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

type DeadShow struct {
	Identifier *string
	Server     *string
	Directory  *string
	Details    DeadShowDetails
	Reviews    []DeadShowReview
	Files      []DeadShowFile
}

type DeadShowDetails struct {
	Title           *string
	Creator         *string
	Mediatype       []string
	Collection      []string
	Type            []string
	Description     *string
	Date            *time.Time
	Year            *uint64
	Subject         *string
	PublicDate      *time.Time
	AddedDate       *time.Time
	Venue           *string
	Coverage        *string
	Source          *string
	Lineage         *string
	Taper           *string
	Transferer      *string
	RunTime         *time.Duration
	Notes           []string
	UpdateDate      *time.Time
	Updater         *string
	Curation        []string
	BackupLocation  *string
	ReviewCount     *uint32
	AverageReview   *float64
	TotalDownloads  *uint32
	DownloadsWeek   *uint32
	DownloadsMonth  *uint32
	ImageLink       *string
	CollectionTitle *string
}

type DeadShowReview struct {
	Title    *string
	Reviewer *string
	Date     *time.Time
	Stars    *uint64
	Body     *string
}

type DeadShowFile struct {
	Name     *string
	Source   *string
	Creator  *string
	Title    *string
	Track    *uint64
	Album    *string
	Bitrate  *uint64
	Format   *string
	Original *string
	MTime    *uint64
	Size     *uint64
	MD5      *string
	CRC32    *string
	SHA1     *string
	Length   *time.Duration
	Height   *uint64
	Width    *uint64
}
