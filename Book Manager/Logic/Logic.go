package Logic

import (
	"Book_Manager/Database"
	"fmt"
	"time"
)

type Books struct {
	Books []Book `json:"books" binding:"required"`
}
type Book struct {
	Name        string    `json:"name" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Volume      int       `json:"volume" binding:"required"`
	PublishedAt time.Time `json:"published_at" binding:"required"`
	Summary     string    `json:"summary" binding:"required"`
	Publisher   string    `json:"publisher" binding:"required"`
}

func SerializeArrayToString(arr []string) string {
	str := ""
	for _, s := range arr {
		str += s + ","
	}
	return str
}
func UnSerializeArrayToString(str string) *[]string {
	var arr []string
	temp := ""
	for i := 0; i < len(str); i++ {
		if str[i] == ',' {
			arr = append(arr, temp)
			temp = ""
		} else {
			temp += string(str[i])
		}
	}
	return &arr

}
func CreateResponseBooks(u Database.UserDb) Books {
	var books Books
	for _, book := range u.Books {
		books.Books = append(books.Books, Book{
			Name:        book.Name,
			Author:      fmt.Sprintf("%s %s", book.Author.Firstname, book.Author.Lastname),
			Category:    book.Category,
			Volume:      book.Volume,
			PublishedAt: book.PublishedAt,
			Summary:     book.Summary,
			Publisher:   book.Publisher,
		})
	}
	return books
}
