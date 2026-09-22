package extraction

import (
	"fmt"
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

func ExtractCover(file multipart.File) ([]byte, error) {

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	defer file.Seek(0, io.SeekStart)

	var coverPDFBytes []byte
	coverPage := []string{"1"}

	digestReader := func(reader io.Reader, page int) error {
		var err error
		// Read all bytes of the 1-page PDF into our slice
		coverPDFBytes, err = io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to read page %d: %w", page, err)
		}
		return nil
	}

	err := api.ExtractPages(file, coverPage, digestReader, nil)

	if err != nil {
		return nil, err
	}

	if len(coverPDFBytes) == 0 {
    return nil, fmt.Errorf("no cover page content could be extracted")
 }

	return coverPDFBytes, nil
}


