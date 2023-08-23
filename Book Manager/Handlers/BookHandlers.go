package Handlers

import (
	"Book_Manager/Database"
	"Book_Manager/Logic"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// CREATE REQUEST BODY

type Book struct {
	Name            string `json:"name" binding:"required"`
	Author          `json:"author" binding:"required"`
	Category        string    `json:"category" binding:"required"`
	Volume          int       `json:"volume" binding:"required"`
	PublishedAt     time.Time `json:"published_at" binding:"required"`
	Summary         string    `json:"summary" binding:"required"`
	TableOfContents []string  `json:"table_of_contents" binding:"required"`
	Publisher       string    `json:"publisher" binding:"required"`
}
type Author struct {
	Firstname   string    `json:"first_name" binding:"required"`
	Lastname    string    `json:"last_name" binding:"required"`
	Birthday    time.Time `json:"birthday" binding:"required"`
	Nationality string    `json:"nationality" binding:"required"`
}

// UPDATE REQUEST BODY

type UpdateBookRequestBody struct {
	Name             string `json:"name"`
	UpdateBookAuthor `json:"author"`
	Category         string    `json:"category"`
	Volume           int       `json:"volume"`
	PublishedAt      time.Time `json:"published_at"`
	Summary          string    `json:"summary" `
	TableOfContents  []string  `json:"table_of_contents" `
	Publisher        string    `json:"publisher" `
}
type UpdateBookAuthor struct {
	Firstname   string    `json:"first_name" `
	Lastname    string    `json:"last_name" `
	Birthday    time.Time `json:"birthday" `
	Nationality string    `json:"nationality" `
}

func (s *BookMangerServer) CreateBookHandle(c *gin.Context) {
	var reqData Book

	// Authenticate of user
	token := c.GetHeader("access_token")
	username, err := s.Authenticate.AuthenticateUserWithToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "access denied"})
		return
	}

	// unmarshal json
	err = c.BindJSON(&reqData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	// create new user db
	author := Database.Author{

		Firstname:   reqData.Firstname,
		Lastname:    reqData.Lastname,
		Birthday:    reqData.Birthday,
		Nationality: reqData.Nationality,
	}
	book := &Database.Book{

		Name:            reqData.Name,
		Category:        reqData.Category,
		Volume:          reqData.Volume,
		PublishedAt:     reqData.PublishedAt,
		Summary:         reqData.Summary,
		TableOfContents: Logic.SerializeArrayToString(reqData.TableOfContents),
		Publisher:       reqData.Publisher,
		Author:          author,
	}

	//add book to books of user
	err = s.Db.CreateNewBook(book, *username)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusCreated, "")
	return
}

func (s *BookMangerServer) GetAllBookOfUserHandle(c *gin.Context) {

	// Authenticate of user
	token := c.GetHeader("access_token")
	username, err := s.Authenticate.AuthenticateUserWithToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "access denied"})
		return
	}

	// get user with book
	user, err := s.Db.GetUserByUsernameWithBooks(*username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "can not find books"})
		return
	}

	// create response
	res := Logic.CreateResponseBooks(*user)
	c.IndentedJSON(http.StatusOK, res)
	return
}

func (s *BookMangerServer) GetBookByIdHandle(c *gin.Context) {

	// AUTHENTICATE OF USER
	token := c.GetHeader("access_token")
	username, err := s.Authenticate.AuthenticateUserWithToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "access denied"})
		return
	}

	// GET PARAM AND VALIDATE
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is missing"})
		return
	}
	temp, err := strconv.Atoi(id)
	if err != nil || temp < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}
	ID := uint(temp)

	// GET BOOK
	book, err := s.Db.GetBookById(*username, ID)
	if err != nil {
		if err.Error() == "book with this id not fond" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not exist"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "can not find book"})
		return
	}

	// CREATE RESPONSE AND RETURN BOOK
	c.IndentedJSON(http.StatusOK, Book{
		Name: book.Name,
		Author: Author{
			Firstname:   book.Author.Firstname,
			Lastname:    book.Author.Lastname,
			Birthday:    book.Author.Birthday,
			Nationality: book.Author.Nationality,
		},
		Category:        book.Category,
		Volume:          book.Volume,
		PublishedAt:     book.PublishedAt,
		Summary:         book.Summary,
		TableOfContents: *Logic.UnSerializeArrayToString(book.TableOfContents),
		Publisher:       book.Publisher,
	})
	return

}

func (s *BookMangerServer) UpdateBookHandle(c *gin.Context) {
	var reqData UpdateBookRequestBody

	// AUTHENTICATE OF USER
	token := c.GetHeader("access_token")
	username, err := s.Authenticate.AuthenticateUserWithToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "access denied"})
		return
	}

	// GET PARAM AND VALIDATE
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is missing"})
		return
	}
	temp, err := strconv.Atoi(id)
	if err != nil || temp < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}
	ID := uint(temp)

	// unmarshal json
	err = c.BindJSON(&reqData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	err = s.Db.UpdateBookOfUser(*username, ID, Database.Book{

		Name:            reqData.Name,
		Category:        reqData.Category,
		Volume:          reqData.Volume,
		PublishedAt:     reqData.PublishedAt,
		Summary:         reqData.Summary,
		TableOfContents: Logic.SerializeArrayToString(reqData.TableOfContents),
		Publisher:       reqData.Publisher,
		Author: Database.Author{

			Firstname:   reqData.Firstname,
			Lastname:    reqData.Lastname,
			Birthday:    reqData.Birthday,
			Nationality: reqData.Nationality,
		},
	})
	if err != nil {
		if err.Error() == "book with this id not fond" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not exist"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, "unable to update book")
		return
	}
	c.IndentedJSON(http.StatusOK, "")
	return

}

func (s *BookMangerServer) DeleteBookHandle(c *gin.Context) {
	// AUTHENTICATE OF USER
	token := c.GetHeader("access_token")
	username, err := s.Authenticate.AuthenticateUserWithToken(token)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "access denied"})
		return
	}

	// GET PARAM AND VALIDATE
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is missing"})
		return
	}
	temp, err := strconv.Atoi(id)
	if err != nil || temp < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}
	ID := uint(temp)

	if err := s.Db.DeleteBookOfUser(*username, ID); err != nil {

		if err.Error() == "book with this id not fond" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not exist"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "unable to delete book"})
		return
	}
	c.IndentedJSON(http.StatusOK, "")
	return
}
