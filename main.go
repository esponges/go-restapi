package main

import (
	"encoding/json"
	"net/http"
	"restapi/auth"
	"restapi/middleware"
	"strconv"

	"github.com/gorilla/mux"
)

// Book represents a book object with title, author, and ISBN fields
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book
var nextID int = 1

func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	http.NotFound(w, r)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = nextID
	nextID++
	books = append(books, book)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedBook Book
	// check for errors and assign the body to the updatedBook variable
	// Decode(&updatedBook) will decode the JSON body and assign it to the updatedBook variable
	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, book := range books {
		if book.ID == id {
			updatedBook.ID = id
			books[i] = updatedBook
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}

	http.NotFound(w, r)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	// Initialize router
	r := mux.NewRouter()

	// Mock data
	books = append(books, Book{ID: 1, Title: "To Kill a Mockingbird", Author: "Harper Lee"})
	books = append(books, Book{ID: 2, Title: "1984", Author: "George Orwell"})
	books = append(books, Book{ID: 3, Title: "Brave New World", Author: "Aldous Huxley"})

	// Route handlers
	r.HandleFunc("/books", middleware.Logging(getBooks)).Methods("GET")
	// use middleware wrapper
	r.HandleFunc("/books-wrapper", middleware.Wrapper(getBooks, middleware.Method("GET")))
	r.HandleFunc("/books/{id}", middleware.Logging(getBook)).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	// auth routes
	r.HandleFunc("/login", auth.Login)
	r.HandleFunc("/logout", auth.Logout)
	r.HandleFunc("/secret-msg", auth.Secret)
	// hashpw
	r.HandleFunc("/hashpw/{pw}", auth.HashPasswordHandler).Methods("GET")

	// Start server
	http.ListenAndServe(":8000", r)
}
