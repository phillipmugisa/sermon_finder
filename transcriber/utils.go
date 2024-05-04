package transcriber

import (
	"strings"
	"time"
)

type Sermon struct {
	Name   string
	Length int
	rate   time.Duration

	FilePath   string
	FileFormat string

	Transcribed bool
	UploadedOn  time.Time
}

func NewSermon(filepath string, fileformat string, Length int, rate time.Duration, uploadedOn time.Time, transcribed bool) *Sermon {
	return &Sermon{
		Name:        getSermonName(filepath),
		Length:      Length,
		rate:        rate,
		FilePath:    filepath,
		FileFormat:  fileformat,
		UploadedOn:  uploadedOn,
		Transcribed: transcribed,
	}
}

func getSermonName(filepath string) string {
	parts := strings.Split(filepath, "/")
	// remove file extension
	name := strings.Split(parts[len(parts)-1], ".")
	return name[0]
}
