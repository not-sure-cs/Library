package api

import (
	"errors"
	"slices"

	"github.com/knibirdgautam/library/internal/database"
)

func validateBook(b database.Book) bool {

	if b.Title == "" || b.Author == "" {
		return false
	}

	if len(b.Title) > 100 || len(b.Author) > 200 {

		return false
	}

	return true
}

func findBook(db *[]database.Book, id int) (*database.Book, error) {

	for i := range *db {
		if id == (*db)[i].ID {
			return &(*db)[i], nil
		}
	}

	return nil, errors.New("Book not found")
}

func delBook(db *[]database.Book, id int) error {

	for i, element := range *db {
		if id == element.ID {
			*db = slices.Delete(*db, i, i+1)
			return nil
		}
	}

	return errors.New("Book not found")
}
