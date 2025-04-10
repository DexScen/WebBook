package service

import (
	"context"

	"github.com/DexScen/WebBook/internal/domain"
)

type BooksRepository interface{
	GetBooks(ctx context.Context, list *domain.ListBooks) error
	PostBook(ctx context.Context, book *domain.Book) error
	DeleteBookByID(ctx context.Context, id int) error
	PatchBook(ctx context.Context, book *domain.Book) error
}

type Books struct{
	repo BooksRepository
}

func NewBooks(repo BooksRepository) *Books{
	return &Books{
		repo: repo,
	}
}

func (b *Books) GetBooks (ctx context.Context, list *domain.ListBooks) error{
	return b.repo.GetBooks(ctx, list)
}

func (b *Books) PostBook (ctx context.Context, book *domain.Book) error{
	return b.repo.PostBook(ctx, book)
}

func (b *Books) DeleteBookByID (ctx context.Context, id int) error{
	return b.repo.DeleteBookByID(ctx, id)
}

func (b *Books) PatchBook (ctx context.Context, book *domain.Book) error{
	return b.repo.PatchBook(ctx, book)
}