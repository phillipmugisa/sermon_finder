package analyzer

import (
	"strings"
	"time"
)

const (
	windowSize       = 4096 // Size of FFT window (must be a power of 2)
	hopSize          = 2048 // Step size for sliding the FFT window
	peakThreshold    = 0.8  // Minimum relative amplitude to consider as a peak
	constellationFan = 10   // Number of peaks to pair for fingerprint
)

type Peak struct {
	FrequencyBin int     // Frequency bin index
	Time         int     // Time (frame index)
	Magnitude    float64 // Amplitude at this peak
}

type Fingerprint struct {
	Hash uint64 // Generated hash from peak pairs
	Time int    // Anchor time of the fingerprint
}

type Sermon struct {
	Name   string
	Length int
	rate   time.Duration

	FilePath   string
	FileFormat string

	Analyzed   bool
	UploadedOn time.Time
}

func NewSermon(filepath string, fileformat string, Length int, rate time.Duration, uploadedOn time.Time, analyzed bool) *Sermon {
	return &Sermon{
		Name:       getSermonName(filepath),
		Length:     Length,
		rate:       rate,
		FilePath:   filepath,
		FileFormat: fileformat,
		UploadedOn: uploadedOn,
		Analyzed:   analyzed,
	}
}

func getSermonName(filepath string) string {
	parts := strings.Split(filepath, "/")
	// remove file extension
	name := strings.Split(parts[len(parts)-1], ".")
	return name[0]
}

type AudioProcessingError struct {
	msg string
}

func NewAudioProcessingError(msg string) *AudioProcessingError {
	return &AudioProcessingError{
		msg: msg,
	}
}

func (err *AudioProcessingError) Error() string {
	return err.msg
}
