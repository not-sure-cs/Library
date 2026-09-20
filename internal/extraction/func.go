package extraction

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type BookMetaData struct {
	Title      string
	Author     string
	Producer   string
	Subject    string
	PageCount  int
	PDFVersion string
}

func ExtractMime(file multipart.File) (string, error) {

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {

		return "", err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buffer[:n])
	return mimeType, nil
}

func ExtractMetadata(file multipart.File, fileHandler *multipart.FileHeader) (*BookMetaData, error) {

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	defer file.Seek(0, io.SeekStart)

	fileName := ""
	if fileHandler != nil {
		fileName = fileHandler.Filename
	}

	info, err := api.PDFInfo(file, fileName, nil, false, nil)
	if err != nil {
		return nil, err
	}

	metadata := &BookMetaData{
		Title:      info.Title,
		Author:     info.Author,
		Subject:    info.Subject,
		Producer:   info.Producer,
		PageCount:  info.PageCount,
		PDFVersion: info.Version,
	}

	return metadata, nil

}
