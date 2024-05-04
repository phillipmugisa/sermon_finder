package transcriber

import (
	"fmt"
	"os"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func Transcribe(file *os.File) {
	// use id to fetch sermon data from database
	// filepath to access file

	file, err := os.Open(file.Name())
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a decoder for the WAV format
	decoder := wav.NewDecoder(file)
	if decoder == nil {
		fmt.Println("Error creating decoder")
		return
	}

	// Get the audio format information
	format := decoder.Format()
	if format == nil {
		fmt.Println("Error getting format")
		return
	}

	// Calculate total number of samples
	totalSamples := decoder.NumChans

	// Define chunk duration (2 minutes)
	chunkDuration := 2 * time.Minute
	chunkSamples := uint16(chunkDuration.Seconds() * float64(format.SampleRate))

	fmt.Printf("totalSamples: %d\n\n", totalSamples)

	// Process audio in 2-minute chunks
	for startSample := uint16(0); startSample < totalSamples; startSample += chunkSamples {
		// Read audio chunk
		audioBuffer := new(audio.IntBuffer)
		_, err := decoder.PCMBuffer(audioBuffer)
		if err != nil {
			fmt.Println("Error reading audio chunk:", err)
			return
		}
		fmt.Printf("audioBuffer: %+v\n\n", audioBuffer)
	}
}
