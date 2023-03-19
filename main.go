package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
    ID     string  `json:"id,omitempty"`
    Title  string  `json:"title,omitempty"`
    Author string  `json:"author,omitempty"`
    ISBN   string  `json:"isbn,omitempty"`
}

var books []Book

func main() {
    // Initialize router
    router := mux.NewRouter()

    // Add some sample data
    books = append(books, Book{ID: "1", Title: "Book One", Author: "Author One", ISBN: "12345"})
    books = append(books, Book{ID: "2", Title: "Book Two", Author: "Author Two", ISBN: "67890"})

    // Define routes
    router.HandleFunc("/books", getBooks).Methods("GET")
    router.HandleFunc("/books/{id}", getBook).Methods("GET")
    router.HandleFunc("/books", createBook).Methods("POST")
    router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
    router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

    // Start server
    log.Fatal(http.ListenAndServe(":8000", router))
}

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(books)
}

// Get single book by ID
func getBook(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getBook() called")

    params := mux.Vars(r)
    for _, book := range books {
        if book.ID == params["id"] {
            json.NewEncoder(w).Encode(book)
            return
        }
    }
    json.NewEncoder(w).Encode(nil)
}

// Create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
    var book Book
    json.NewDecoder(r.Body).Decode(&book)
    books = append(books, book)
    json.NewEncoder(w).Encode(book)
}

// Update an existing book by ID
func updateBook(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for i, book := range books {
        if book.ID == params["id"] {
            books[i] = Book{
                ID:     book.ID,
                Title:  book.Title,
                Author: book.Author,
                ISBN:   book.ISBN,
            }
            json.NewDecoder(r.Body).Decode(&books[i])
            json.NewEncoder(w).Encode(books[i])
            return
        }
    }
    json.NewEncoder(w).Encode(nil)
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
