package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/DexScen/WebBook/backend/books/internal/domain"
	"github.com/gorilla/mux"
)

type Books interface {
	GetBooks(ctx context.Context, list *domain.ListBooks) error
	PostBook(ctx context.Context, book *domain.Book) error
	DeleteBookByID(ctx context.Context, id int) error
	PutBook(ctx context.Context, book *domain.Book) error
}

type Handler struct {
	booksService Books
}

func NewHandler(books Books) *Handler {
	return &Handler{
		booksService: books,
	}
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (h *Handler) OptionsHandler(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Use(loggingMiddleware)

	links := r.PathPrefix("/books").Subrouter()
	{
		links.HandleFunc("", h.GetBooks).Methods(http.MethodGet)
		links.HandleFunc("", h.PostBook).Methods(http.MethodPost)
		links.HandleFunc("", h.DeleteBookByID).Methods(http.MethodDelete)
		links.HandleFunc("", h.PutBook).Methods(http.MethodPut)
		links.HandleFunc("", h.OptionsHandler).Methods(http.MethodOptions)
	}
	return r
}

func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var list domain.ListBooks
	if err := h.booksService.GetBooks(context.TODO(), &list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getBooks error:", err)
		return
	}

	if jsonResp, err := json.Marshal(list); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("getBooks error:", err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func (h *Handler) PostBook(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var book domain.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("postBook error:", err)
		return
	}

	if err := h.booksService.PostBook(context.TODO(), &book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("postBook error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteBookByID(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var req domain.DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("deleteBookByID error:", err)
		return
	}

	if err := h.booksService.DeleteBookByID(context.TODO(), req.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("deleteBookByID error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PutBook(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	var book domain.Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("PatchBook error:", err)
		return
	}

	if err := h.booksService.PutBook(context.TODO(), &book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("PatchBook error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
