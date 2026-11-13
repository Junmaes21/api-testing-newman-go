package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type envelope struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type bookStore struct {
	mu     sync.RWMutex
	books  map[int]Book
	nextID int
}

func NewHandler() http.Handler {
	store := &bookStore{
		books: map[int]Book{
			1: {ID: 1, Title: "La ciudad y los perros", Author: "Mario Vargas Llosa"},
			2: {ID: 2, Title: "El principito", Author: "Antoine de Saint-Exupéry"},
		},
		nextID: 3,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", serveHome)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /api/books", store.listBooks)
	mux.HandleFunc("POST /api/books", store.createBook)
	mux.HandleFunc("GET /api/books/{id}", store.getBook)
	mux.HandleFunc("PUT /api/books/{id}", store.updateBook)
	mux.HandleFunc("DELETE /api/books/{id}", store.deleteBook)
	return mux
}

func (s *bookStore) listBooks(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	books := make([]Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	writeJSON(w, http.StatusOK, envelope{Data: books})
}

func (s *bookStore) getBook(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	s.mu.RLock()
	book, found := s.books[id]
	s.mu.RUnlock()
	if !found {
		writeError(w, http.StatusNotFound, "book not found")
		return
	}
	writeJSON(w, http.StatusOK, envelope{Data: book})
}

func (s *bookStore) createBook(w http.ResponseWriter, r *http.Request) {
	book, ok := decodeBook(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	book.ID = s.nextID
	s.nextID++
	s.books[book.ID] = book
	s.mu.Unlock()
	w.Header().Set("Location", "/api/books/"+strconv.Itoa(book.ID))
	writeJSON(w, http.StatusCreated, envelope{Data: book})
}

func (s *bookStore) updateBook(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	book, ok := decodeBook(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	if _, found := s.books[id]; !found {
		s.mu.Unlock()
		writeError(w, http.StatusNotFound, "book not found")
		return
	}
	book.ID = id
	s.books[id] = book
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, envelope{Data: book})
}

func (s *bookStore) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	if _, found := s.books[id]; !found {
		s.mu.Unlock()
		writeError(w, http.StatusNotFound, "book not found")
		return
	}
	delete(s.books, id)
	s.mu.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}
	return id, true
}

func decodeBook(w http.ResponseWriter, r *http.Request) (Book, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()
	var book Book
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		writeError(w, http.StatusBadRequest, "request body must be valid JSON")
		return Book{}, false
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "request body must contain one JSON object")
		return Book{}, false
	}
	book.Title = strings.TrimSpace(book.Title)
	book.Author = strings.TrimSpace(book.Author)
	if book.Title == "" || book.Author == "" {
		writeError(w, http.StatusBadRequest, "title and author are required")
		return Book{}, false
	}
	return book, true
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, envelope{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
