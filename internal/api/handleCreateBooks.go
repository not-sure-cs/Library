package api

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/knibirdgautam/library/internal/database"
	"github.com/knibirdgautam/library/internal/extraction"
	"github.com/knibirdgautam/library/internal/storage"
)

func HandleCreateBooks(queries database.DBQueries, store storage.R2Store, secret storage.Secret) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RespondWithError(w, http.StatusMethodNotAllowed, "Only POST requests allowed")
			return
		}

		err := r.ParseMultipartForm(200 << 20)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed to parse multipart form")
			return
		}

		file, fileHandler, err := r.FormFile("uploadFile")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "uploadFile is required")
			return
		}
		defer file.Close()

		mimeType, err := extraction.ExtractMime(file)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Mime data couldn't be extracted")
			return
		}

		fileHash, err := extraction.GenerateSHA256FileHash(fileHandler)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed the Filehash function")
			return
		}

		exists, err := queries.CheckApiKeyExists(r.Context(), fileHash)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Failed the key check")
		}

		if exists != false {
			log.Print("File Already Exists")
			return
		}

		var meta *extraction.BookMetaData
		if mimeType == "application/pdf" {
			meta, err = extraction.ExtractMetadata(file, fileHandler)
			if err != nil {
				log.Printf("Failed to extract PDF metadata: %v", err)
			}
		}

		title := ""
		authorName := ""
		if meta != nil {
			title = strings.TrimSpace(meta.Title)
			authorName = strings.TrimSpace(meta.Author)
		}

		if title == "" && fileHandler != nil {
			filename := fileHandler.Filename
			ext := filepath.Ext(filename)
			title = strings.TrimSuffix(filename, ext)
		}
		if title == "" {
			title = "Untitled"
		}

		if authorName == "" {
			authorName = "Unknown"
		}

		if _, err := file.Seek(0, io.SeekStart); err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to rewind file")
			return
		}

		fileKey, err := database.SaveFile(0, r.Context(), secret, store, file, fileHandler)
		if err != nil {
			log.Printf("R2 Upload Error: %v", err)
			RespondWithError(w, http.StatusInternalServerError, "Upload failed")
			return
		}

		author, err := queries.GetAuthor(r.Context(), authorName)
		if err != nil {
			author, err = queries.CreateAuthor(r.Context(), database.CreateAuthorParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      authorName,
			})

			if err != nil {
				_ = database.UnsaveFile(r.Context(), secret, store, fileKey)
				RespondWithError(w, http.StatusInternalServerError, "Could not Create Author")
				return
			}
		}

		var pageCount sql.NullInt32
		var producer, subject, pdfVersion sql.NullString
		if meta != nil {
			if meta.PageCount > 0 {
				pageCount = database.ToNullInt32(int32(meta.PageCount))
			}
			producer = database.ToNullString(meta.Producer)
			subject = database.ToNullString(meta.Subject)
			pdfVersion = database.ToNullString(meta.PDFVersion)
		}

		book, err := queries.CreateBook(r.Context(), database.CreateBookParams{
			ID:         uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Name:       title,
			Isbn:       database.ToNullString(""),
			FilePath:   fileKey,
			MimeType:   database.ToNullString(mimeType),
			PageCount:  pageCount,
			Producer:   producer,
			Subject:    subject,
			PdfVersion: pdfVersion,
		})

		if err != nil {
			_ = database.UnsaveFile(r.Context(), secret, store, fileKey)
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("DBError: %s", err))
			return
		}

		linker, err := queries.LinkBookAuthor(r.Context(), database.LinkBookAuthorParams{
			BookID:   book.ID,
			AuthorID: author.ID,
			ApiKey:   fileHash,
		})

		if err != nil {
			_ = database.UnsaveFile(r.Context(), secret, store, fileKey)
			RespondWithError(w, http.StatusInternalServerError, "Couldn't link Books and Authors")
			return
		}

		resp := database.Linked{
			Author:   author,
			Book:     book,
			Link:     linker,
			Metadata: meta,
		}

		RespondWithJSON(w, http.StatusCreated, resp)
	}
}
