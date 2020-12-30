package main

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "sync"
    "time"
)

func StartWebServer(port int, wg *sync.WaitGroup, errChnl chan <- error) (*http.Server, error) {

    route := mux.NewRouter()

    route.HandleFunc("/authors", HandleListAuthors).Methods("GET")
    route.HandleFunc("/authors/{id}", HandleGetAuthor).Methods("GET")
    route.HandleFunc("/books", HandleListBooks).Methods("GET")
    route.HandleFunc("/books/{id}", HandleGetBook).Methods("GET")
    route.HandleFunc("/books-by-author/{id}", HandleGetBooksByAuthor).Methods("GET")
    route.HandleFunc("/books/add", HandleAddNewBook).Methods("POST")
    route.HandleFunc("/authors/add", HandleAddNewAuthor).Methods("POST")

    //handler := cors.Default().Handler(route)

    server := &http.Server{
        //Handler: handler,
        Handler: route,
        Addr:         fmt.Sprintf("localhost:%d", port),
        WriteTimeout: 30 * time.Second,
        ReadTimeout:  30 * time.Second,
    }

    log.Printf("Starting Webserver....")
    wg.Add(1)
    go runServer(server, errChnl, wg)

    return server, nil
}

func runServer(server *http.Server, errChnl chan <- error, wg *sync.WaitGroup) {
    log.Printf("Calling ListenAndServe()!")
    errChnl <- server.ListenAndServe()
    wg.Done()
}

func HandleListBooks(w http.ResponseWriter, _ *http.Request) {
    books, err := ReadTableBooks()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(fmt.Sprintf("Unable to open the PG-DB!! Error - %v", err)))
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(&books)
}

func HandleListAuthors(w http.ResponseWriter, _ *http.Request) {
    authors, err := ReadTableAuthors()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(fmt.Sprintf("Unable to open the PG-DB!! Error - %v", err)))
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(authors)
}

func HandleGetAuthor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, ok := params["id"]
    msg := ""
    if ok {
        author, err := GetAuthorById(id)
        if err == nil {
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(author)
            return
        }
        msg = fmt.Sprintf("Error - %v", err)
        w.WriteHeader(http.StatusNotFound)
    } else {
        msg = fmt.Sprintf("id parameter is required!")
        w.WriteHeader(http.StatusBadRequest)
    }
    w.Write([]byte(msg))
    return
}

func HandleGetBook(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, ok := params["id"]
    msg := ""
    if ok {
        book, err := GetBookById(id)
        if err == nil {
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(book)
            return
        }
        msg = fmt.Sprintf("Error - %v", err)
        w.WriteHeader(http.StatusNotFound)
    } else {
        msg = fmt.Sprintf("id parameter is required!")
        w.WriteHeader(http.StatusBadRequest)
    }
    w.Write([]byte(msg))
    return
}

func HandleGetBooksByAuthor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, ok := params["id"]
    msg := ""
    if ok {
        authors, err := GetBooksByAuthor(id)
        if err == nil {
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(authors)
            return
        }
        msg = fmt.Sprintf("%v", err)
        w.WriteHeader(http.StatusNotFound)
    } else {
        msg = fmt.Sprintf("id parameter is required!")
        w.WriteHeader(http.StatusBadRequest)
    }
    w.Write([]byte(msg))
    return
}

func HandleAddNewAuthor(w http.ResponseWriter, r *http.Request) {

    var retNow bool
    var msg string
    if r.Method != "POST" {
        msg := fmt.Sprintf("/add url requires 'POST' operation!")
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(msg))
        return
    }

    newAuthor := &Author{}
    err := json.NewDecoder(r.Body).Decode(newAuthor)
    if err != nil {
        msg = fmt.Sprintf("Error - %v", err)
        retNow = true
    }

    if newAuthor.ID != 0 {
        msg = fmt.Sprintf("Error - Author ID primary key cannot be non-zero!")
        retNow = true
    }
    if newAuthor.FirstName == "" || newAuthor.LastName == "" {
        msg = fmt.Sprintf("Error - Author's First and Last Name are required!")
        retNow = true
    }
    if retNow {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(msg))
        return
    }

    newAuthor.CreatedAt = time.Time{}
    newAuthor.UpdatedAt = time.Time{}
    newAuthor.DeletedAt = nil

    author, err := AddNewAuthor(newAuthor)
    if err != nil {
        msg = fmt.Sprintf("Error - DB Write Failed :: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(author)
    return
}

func HandleAddNewBook (w http.ResponseWriter, r *http.Request) {

    var retNow bool
    var msg string
    if r.Method != "POST" {
        msg := fmt.Sprintf("/add url requires 'POST' operation!")
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(msg))
        return
    }

    newBook := &Book{}
    err := json.NewDecoder(r.Body).Decode(newBook)
    if err != nil {
        msg = fmt.Sprintf("Error - %v", err)
        retNow = true
    }

    if newBook.ID != 0 {
        msg = fmt.Sprintf("Error - Book ID primary key cannot be non-zero!")
        retNow = true
    }

    if newBook.Title == "" {
        msg = fmt.Sprintf("Error - Book Title is required!")
        retNow = true
    }

    if len(newBook.Authors) == 0 || len(newBook.Authors) > 10 {
        msg = fmt.Sprintf("Error - Book Authors not provided or list is too large(max 10 authors)!")
        retNow = true
    }
    for _, authorId := range newBook.Authors {
        authId := fmt.Sprintf("%d", authorId)
        _, err := GetAuthorById(authId)
        if err != nil {
            msg = fmt.Sprintf("Error - AuthorId %d uknown! Pl. create the Author Identify first!", authorId)
            retNow = true
        }
    }
    if retNow {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(msg))
        return
    }

    newBook.CreatedAt = time.Time{}
    newBook.UpdatedAt = time.Time{}
    newBook.DeletedAt = nil

    book, err := AddNewBook(newBook)
    if err != nil {
        msg = fmt.Sprintf("Error - DB Write Failed :: %v", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(book)
    return
}