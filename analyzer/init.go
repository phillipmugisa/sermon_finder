package analyzer

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/go-mp3"
	"github.com/mjibson/go-dsp/fft"
	"github.com/xuri/excelize/v2"
)

func toExcel(spectrogram [][]float64, peaks []Peak) {

	// Create a new Excel file
	f := excelize.NewFile()

	// Add spectrogram data to "Spectrogram" sheet
	// sheetSpectrogram := "Spectrogram"
	// f.NewSheet(sheetSpectrogram)
	// for i, row := range spectrogram {
	// 	for j, value := range row {
	// 		cellName, err := excelize.ColumnNumberToName(j)
	// 		if err != nil {
	// 			cellName = "tt"
	// 		}
	// 		cell := fmt.Sprintf("%s%d", cellName, i+1)
	// 		f.SetCellValue(sheetSpectrogram, cell, value)
	// 	}
	// }

	// Add peaks data to "Peaks" sheet
	sheetPeaks := "Peaks"
	f.NewSheet(sheetPeaks)
	// Add headers
	headers := []string{"Time", "Frequency", "Magnitude"}
	for i, header := range headers {
		cellName, err := excelize.ColumnNumberToName(i)
		if err != nil {
			cellName = "tt"
		}
		cell := fmt.Sprintf("%s1", cellName)
		f.SetCellValue(sheetPeaks, cell, header)
	}
	// Add peak values
	for i, peak := range peaks {
		f.SetCellValue(sheetPeaks, fmt.Sprintf("A%d", i+2), peak.Time)
		f.SetCellValue(sheetPeaks, fmt.Sprintf("B%d", i+2), peak.FrequencyBin)
		f.SetCellValue(sheetPeaks, fmt.Sprintf("C%d", i+2), peak.Magnitude)
	}

	// Save the file
	if err := f.SaveAs("spectrogram_and_peaks.xlsx"); err != nil {
		log.Fatalf("failed to save Excel file: %v", err)
	}
}

func AnalyzeAudio(filePath string) error {
	// use id to fetch sermon data from database
	// filepath to access file
	file, err := os.Open(filePath)
	if err != nil {
		return NewAudioProcessingError(fmt.Sprintf("Unable to access audio file %v.", filePath))
	}
	defer file.Close()

	samples, _, err := getFileDataSamples(file)
	if err != nil {
		return NewAudioProcessingError(fmt.Sprintf("Unable to process audio file %v.", filePath))
	}

	spectrogram, computeErr := computeSpectrogram(samples)
	if computeErr != nil {
		return NewAudioProcessingError(fmt.Sprintf("Unable to generate spectrogram for  %v.", filePath))
	}
	fmt.Println("spectrogram generated")

	peaks, computePeaksErr := findSpectrogramPeaks(spectrogram)
	if computePeaksErr != nil {
		return NewAudioProcessingError(fmt.Sprintf("Unable to generate peaks for  %v.", filePath))
	}
	fmt.Println("peaks obtained")

	fmt.Println("exceling")
	toExcel(spectrogram, peaks)

	// fingerprints, generateFPrintErr := generateFingerprints(peaks)
	// if generateFPrintErr != nil {
	// 	return NewAudioProcessingError(fmt.Sprintf("Unable to generate fingerprint for  %v.", filePath))
	// }

	return nil
}

func getFileDataSamples(file *os.File) ([]float32, int, error) {
	// Decode the MP3 file
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		fmt.Println("Error creating MP3 decoder:", err)
		return nil, 0, err
	}
	pcmBuffer := make([]byte, 4096)
	var samples []float32
	for {
		n, err := decoder.Read(pcmBuffer)
		if n == 0 || err != nil {
			break
		}

		for _, data := range pcmBuffer {
			value := int16(data)
			samples = append(samples, float32(value))
		}
	}

	// MP3 typically has a default sample rate of 44100 Hz
	sampleRate := 44100
	return samples, sampleRate, nil
}

func computeSpectrogram(samples []float32) ([][]float64, error) {
	numFrames := (len(samples) - windowSize) / hopSize
	spectrogram := make([][]float64, numFrames)

	for i := 0; i < numFrames; i++ {
		start := i * hopSize
		end := start + windowSize
		window := samples[start:end]

		output := make([]float64, len(window))
		for i, value := range window {
			output[i] = float64(value)
		}

		fftResult := fft.FFTReal(output)

		// Convert complex FFT output to magnitude
		// we only use the first half because output is symmetric for real inputs
		magnitude := make([]float64, len(fftResult)/2)
		for j := 0; j < len(magnitude); j++ {
			// we get the magnitude of each = square root of (i^2 + j^2)
			real, imag := real(fftResult[j]), imag(fftResult[j])
			magnitude[j] = math.Sqrt(real*real + imag*imag)
		}
		spectrogram[i] = magnitude
	}

	return spectrogram, nil
}

// findSpectrogramPeaks identifies prominent peaks in the spectrogram
func findSpectrogramPeaks(spectrogram [][]float64) ([]Peak, error) {
	var peaks []Peak
	for t, frame := range spectrogram {
		for f, magnitude := range frame {
			// Check if magnitude exceeds threshold and is a local maximum
			if magnitude >= peakThreshold && isLocalMaximum(spectrogram, t, f) {
				peaks = append(peaks, Peak{FrequencyBin: f, Time: t, Magnitude: magnitude})
			}
		}
	}
	return peaks, nil
}

/*

	-1| 0 | 1

-1 	* | * | *
   	----------
0	* | * | *
	----------
+1	* | * | *

dt: -1 to +1 (previous, current, next time frame) <- horizontal
df: -1 to +1 (lower, current, higher frequency bins) <- vertical
*/

// isLocalMaximum checks if a magnitude is a local maximum in a small neighborhood
func isLocalMaximum(spectrogram [][]float64, t, f int) bool {
	for dt := -1; dt <= 1; dt++ {
		for df := -1; df <= 1; df++ {
			if dt == 0 && df == 0 {
				// skip current point
				continue
			}
			nt, nf := t+dt, f+df
			if nt >= 0 && nt < len(spectrogram) && nf >= 0 && nf < len(spectrogram[nt]) {
				if spectrogram[nt][nf] > spectrogram[t][f] {
					return false
				}
			}
		}
	}
	return true
}

// generateFingerprints creates fingerprints from peaks using constellation method
func generateFingerprints(peaks []Peak) ([]Fingerprint, error) {
	var fingerprints []Fingerprint
	for i, anchor := range peaks {
		// Pair anchor peak with subsequent peaks
		for j := 1; j <= constellationFan && i+j < len(peaks); j++ {
			target := peaks[i+j]
			deltaTime := target.Time - anchor.Time
			hash := generateHash(anchor.FrequencyBin, target.FrequencyBin, deltaTime)
			fingerprints = append(fingerprints, Fingerprint{Hash: hash, Time: anchor.Time})
		}
	}
	return fingerprints, nil
}

// generateHash creates a unique hash for a pair of peaks
/*
the hash contains:
	- anchor freq
	- target freq
	- delta time (relative time ie anchor to target)
*/
func generateHash(f1, f2, deltaTime int) uint64 {
	return uint64(f1)<<32 | uint64(f2)<<16 | uint64(deltaTime)
}
