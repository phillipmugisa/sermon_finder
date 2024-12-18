package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/phillipmugisa/sermon_finder/analyzer"
)

func (server *Server) HomeHandler(c context.Context, w http.ResponseWriter, r *http.Request) error {

	analyzer.AnalyzeAudio("/home/mugisa/sermon_finder/media/sermons/The Mystery Of God's will.mp3")

	err := server.Render(c, w, "index", nil)
	if err != nil {
		return NewError("error rendering template", http.StatusInternalServerError)
	}
	return nil
}

func (server *Server) SermonUploadHandler(c context.Context, w http.ResponseWriter, r *http.Request) error {

	r.ParseMultipartForm(10 << 20)
	// read file
	// create destination
	// write file

	file, handler, filesError := r.FormFile("sermon")
	// check is files were uploaded
	if !errors.Is(filesError, http.ErrMissingFile) {
		if filesError != nil {
			return NewError("error parsing sermon uploaded", http.StatusInternalServerError)
		}
		defer file.Close()

		// save to disk
		save_file_path := filepath.Join("media/sermons", handler.Filename)
		destination, fileCreateErr := os.Create(save_file_path)
		if fileCreateErr != nil {
			return NewError("error creating destination file", http.StatusInternalServerError)
		}
		defer destination.Close()

		// Copy the file contents to the destination file
		_, fileWriteErr := io.Copy(destination, file)
		if fileWriteErr != nil {
			return NewError("error saving file", http.StatusInternalServerError)
		}

		// analyze in background
		go analyzer.AnalyzeAudio(save_file_path)
	}

	data := map[string]string{
		"message": "Upload Successful",
		"status":  "success",
	}

	return server.Render(c, w, "form-status", data)
}
