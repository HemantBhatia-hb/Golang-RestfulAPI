package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"math/rand"
	"strconv"

	"github.com/gorilla/mux"
)

// Books Struct Models

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`

	// Author struct..
}
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Book

// Get All Books.....
func getAllBooks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get Single Book........

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	params := mux.Vars(r) //Get params

	// loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create Book.....
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) // convert  book's id from integer to string
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

// Update Book....
func updateBook(w http.ResponseWriter, r *http.Request) {

	// combination of deleteBooks and createBooks... if we delete an existing item first and after that recreate it on same place with updated values.
	// it automatically looks like upadte.

	// so first delete pre-existing book....

	w.Header().Set("content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)

			// and after deleting, create a new book.....   Automatically we update the Book

			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]

			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return

		}
	}
	json.NewEncoder(w).Encode(books)
}

// Delete Book.....
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...) // let us suppose i want to delete book whose index is 2
			//   so it append all the books from index o to 1 and from index 3 to lastindex.
			//so automatically 2 indexed will be deleted.

			break
		}

	}

	json.NewEncoder(w).Encode(books)
}
func main() {
	fmt.Println("Starting the Application")

	// Init Router..
	r := mux.NewRouter()
	// Mock data
	books = append(books, Book{ID: "1", Isbn: "44784", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "44785", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})
	r.HandleFunc("/api/books", getAllBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8087", r))

}
