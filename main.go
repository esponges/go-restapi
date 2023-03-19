package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Book represents a book object with title, author, and ISBN fields
type Book struct {
	ID     string `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
	ISBN   string `json:"isbn,omitempty"`
}

// books is an in-memory collection of books
var books []Book

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

// Get a book by ID
func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validate required fields
	if book.Title == "" || book.Author == "" || book.ISBN == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	book.ID = fmt.Sprintf("%d", len(books)+1)
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update an existing book
func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validate required fields
	if book.Title == "" || book.Author == "" || book.ISBN == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	for i, b := range books {
		if b.ID == params["id"] {
			book.ID = b.ID
			books[i] = book
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

// Delete a book by ID
func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, book := range books {
		if book.ID == params["id"] {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	// Initialize router
	r := mux.NewRouter()

	// Mock data
	books = append(books, Book{ID: "1", Title: "To Kill a Mockingbird", Author: "Harper Lee", ISBN: "9780061120084"})
	books = append(books, Book{ID: "2", Title: "1984", Author: "George Orwell", ISBN: "9780451524935"})
	books = append(books, Book{ID: "3", Title: "Brave New World", Author: "Aldous Huxley", ISBN: "9780060850524"})

	// Route handlers
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// Start server
	http.ListenAndServe(":8000", r)
}
