package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "github.com/lib/pq"
    "log"
)

type Book struct {
    gorm.Model
    Title string
    Authors pq.Int64Array `gorm:"type:integer[]"`
    NumPages  uint32
    Genre   string
    Edition string
    Publisher string
    ISBN string
}

type Author struct {
    gorm.Model
    FirstName string
    MiddleName string
    LastName string
    Age  uint8
    Country string
}

func CreateBooksTableEntries(pg_db *gorm.DB) {
    books := []Book {
        {Title: "The GO Programming Language", NumPages: 380, Genre: "Computer Science, Programming Language",
            Authors: pq.Int64Array([]int64{1, 2}), Edition: "12 2020", Publisher: "Addison-Wesley", ISBN: "978-0-13-419044-0"},
        {Title: "Emotional Intelligence", NumPages: 354, Genre: "Psychology, Human Behavior", Authors: pq.Int64Array([]int64{3}),
            Edition: "1995", Publisher: "Blooms Bury", ISBN: "978 0 7475 2830 2"},
    }

    for _, book := range books {
        pg_db.Create(&book)
    }
}

func CreateAuthorsTableEntries(pg_db *gorm.DB) {
    authors := []Author {
        {FirstName: "Alan", MiddleName: "A.A", LastName: "Donovan", Age: 0, Country: "USA"},
        {FirstName: "Brian", MiddleName: "W", LastName: "Kernighan", Age: 0, Country: "USA"},
        {FirstName: "Daniel", MiddleName: "", LastName: "Goleman", Age: 0, Country: "USA"},
    }

    for _, author := range authors {
        pg_db.Create(&author)
    }
}

func DbInitialize() {
    var pg_db *gorm.DB
    var err error

    pg_db, err = gorm.Open( "postgres", "host=localhost port=5432 user=prabhanj dbname=prabhanj sslmode=disable password=")
    if err != nil {
        log.Fatalf("Failed to open the DB!! Error - %v", err)
    }
    defer pg_db.Close()

    pg_db.AutoMigrate(Author{})
    pg_db.AutoMigrate(Book{})

    //CreateBooksTableEntries(pg_db)
    //CreateAuthorsTableEntries(pg_db)
}

func ReadTableBooks() (books[]Book, err error){
    var pg_db *gorm.DB

    pg_db, err = gorm.Open( "postgres", "host=localhost port=5432 user=prabhanj dbname=prabhanj sslmode=disable password=")
    if err != nil {
        return nil, err
    }
    defer pg_db.Close()

    pg_db.Find(&books)
    return books, nil
}

func ReadTableAuthors() (authors[]Author, err error){
    var pg_db *gorm.DB

    pg_db, err = gorm.Open( "postgres", "host=localhost port=5432 user=prabhanj dbname=prabhanj sslmode=disable password=")
    if err != nil {
        return nil, err
    }
    defer pg_db.Close()

    pg_db.Find(&authors)
    return authors, nil
}

func GetAuthorById(key string) (author *Author, err error){
    var pg_db *gorm.DB
    author = &Author{}

    pg_db, err = gorm.Open( "postgres", "host=localhost port=5432 user=prabhanj dbname=prabhanj sslmode=disable password=")
    if err != nil {
        return nil, err
    }
    defer pg_db.Close()
    err = pg_db.First(author, key).Error
    return author, err
}

func GetBookById(key string) (book *Book, err error){
    var pg_db *gorm.DB
    book = &Book{}

    pg_db, err = gorm.Open( "postgres", "host=localhost port=5432 user=prabhanj dbname=prabhanj sslmode=disable password=")
    if err != nil {
        return nil, err
    }
    defer pg_db.Close()
    err = pg_db.First(book, key).Error
    return book, err
}

func GetBooksByAuthor(key string) (books []Book, err error){
    var pg_db *gorm.DB
    //books = make([]Book, 0)

    pg_db, err = gorm.Open( "postgres", "host=localhost port=5432 user=prabhanj dbname=prabhanj sslmode=disable password=")
    if err != nil {
        return nil, err
    }
    defer pg_db.Close()
    err = pg_db.Where("? = any (Authors) ", key).Find(&books).Error
    return books, err
}

func AddNewAuthor(author *Author) (auth *Author, err error) {
    var pg_db *gorm.DB
    pg_db, err = gorm.Open( "postgres", "host=localhost port=5432 user=prabhanj dbname=prabhanj sslmode=disable password=")
    if err != nil {
        return nil, err
    }
    defer pg_db.Close()
    tmp_db := pg_db.Save(author)
    if tmp_db.Error != nil {
        return nil, tmp_db.Error
    }
    auth = tmp_db.Value.(*Author)
    return auth, nil
}

func AddNewBook(book *Book) (bk *Book, err error) {
    var pg_db *gorm.DB
    pg_db, err = gorm.Open( "postgres", "host=localhost port=5432 user=prabhanj dbname=prabhanj sslmode=disable password=")
    if err != nil {
        return nil, err
    }
    defer pg_db.Close()
    tmp_db := pg_db.Save(book)
    if tmp_db.Error != nil {
        return nil, tmp_db.Error
    }
    bk = tmp_db.Value.(*Book)
    return bk, nil
}
