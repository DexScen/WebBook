package msql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/DexScen/WebBook/backend/books/internal/domain"
	_ "github.com/lib/pq"
)

type Books struct {
	db *sql.DB
}

func NewBooks(db *sql.DB) *Books {
	return &Books{db: db}
}

func (b *Books) GetBooks(ctx context.Context, list *domain.ListBooks) error {
	rows, err := b.db.Query("SELECT id, title, author, year FROM books")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var book domain.Book

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			return err
		}

		*list = append(*list, book)
	}
	return err
}

func (b *Books) PostBook(ctx context.Context, book *domain.Book) error {
	tr, err := b.db.Begin()
	if err != nil {
		return err
	}
	fmt.Println(*book)
	statement, err := tr.Prepare("INSERT INTO books (title, author, year) VALUES ($1, $2, $3)")
	if err != nil {
		tr.Rollback()
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(book.Title, book.Author, book.Year)
	if err != nil {
		tr.Rollback()
		return err
	}

	return tr.Commit()
}

func (b *Books) DeleteBookByID(ctx context.Context, id int) error {
	tr, err := b.db.Begin()
	if err != nil {
		return err
	}

	statement, err := tr.Prepare("DELETE FROM books WHERE id = $1")
	if err != nil {
		tr.Rollback()
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		tr.Rollback()
		return err
	}

	return tr.Commit()
}

func (b *Books) PatchBook(ctx context.Context, book *domain.Book) error {
	tr, err := b.db.Begin()
	if err != nil {
		return err
	}

	statement, err := tr.Prepare(`
	UPDATE books 
	SET title = $1, author = $2, year =$3 
	WHERE id = $4
	`)
	if err != nil {
		tr.Rollback()
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(book.Title, book.Author, book.Year, book.ID)
	if err != nil {
		tr.Rollback()
		return err
	}

	return tr.Commit()
}
