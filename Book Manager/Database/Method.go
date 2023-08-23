package Database

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// USER METHOD

func (gdb *GormDB) GetUserByUsername(username string) (*UserDb, error) {
	var user UserDb
	err := gdb.db.Where(&UserDb{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (gdb *GormDB) GetUserByEmail(email string) (*UserDb, error) {
	var user UserDb
	err := gdb.db.Where(&UserDb{Email: email}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (gdb *GormDB) GetUserByPhone(phone string) (*UserDb, error) {
	var user UserDb
	err := gdb.db.Where(&UserDb{PhoneNumber: phone}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (gdb *GormDB) CheckUserDuplicate(username string, phone string, email string) error {

	if _, err := gdb.GetUserByEmail(email); err == nil {
		return errors.New("user already exist ")
	}
	if _, err := gdb.GetUserByUsername(username); err == nil {
		return errors.New("user already exist ")
	}
	if _, err := gdb.GetUserByPhone(phone); err == nil {
		return errors.New("user already exist ")
	}
	return nil
}
func (gdb *GormDB) GetUserByUsernameWithBooks(username string) (u *UserDb, err error) {
	var user UserDb
	err = gdb.db.Model(&UserDb{}).Where(&UserDb{Username: username}).Preload("Books").Preload("Books.Author").Find(&user).Error
	return &user, err
}

// BOOK METHOD

func (gdb *GormDB) CreateNewUser(u *UserDb) error {
	if pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0); err != nil {
		return err
	} else {
		u.Password = string(pw)
	}

	return gdb.db.Create(u).Error

}
func (gdb *GormDB) CreateNewBook(b *Book, username string) error {
	user, err := gdb.GetUserByUsername(username)
	if err != nil {
		return errors.New("user not found")
	}

	user.Books = append(user.Books, *b)

	err = gdb.db.Updates(&user).Error
	if err != nil {
		return errors.New("unable to add the book")
	}
	return nil

}
func (gdb *GormDB) GetBookById(username string, id uint) (*Book, error) {
	user, err := gdb.GetUserByUsernameWithBooks(username)
	if err != nil {
		return nil, err
	}
	for _, b := range user.Books {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, errors.New("book with this id not fond")
}
func (gdb *GormDB) UpdateBookOfUser(username string, id uint, NewBook Book) error {
	book, err := gdb.GetBookById(username, id)
	if err != nil {
		return err
	}

	// SET value for book
	if NewBook.Name != "" {
		book.Name = NewBook.Name
	}
	if NewBook.Publisher != "" {
		book.Publisher = NewBook.Publisher
	}
	if NewBook.TableOfContents != "" {
		book.TableOfContents = NewBook.TableOfContents
	}
	if NewBook.Summary != "" {
		book.Summary = NewBook.Summary
	}
	if fmt.Sprintf("%v", NewBook.PublishedAt) != "0001-01-01T00:00:00Z" {
		book.PublishedAt = NewBook.PublishedAt
	}
	if NewBook.TableOfContents != "" {
		book.TableOfContents = NewBook.TableOfContents
	}
	if NewBook.Volume != 0 {
		book.Volume = NewBook.Volume
	}
	if NewBook.Category != "" {
		book.Category = NewBook.Category
	}
	if NewBook.Author.Firstname != "" {
		book.Author.Firstname = NewBook.Author.Firstname
	}
	if NewBook.Author.Lastname != "" {
		book.Author.Lastname = NewBook.Author.Lastname
	}
	if NewBook.Author.Nationality != "" {
		book.Author.Nationality = NewBook.Author.Nationality
	}
	if fmt.Sprintf("%v", NewBook.Author.Birthday) != "0001-01-01T00:00:00Z" {
		book.Author.Birthday = NewBook.Author.Birthday
	}

	err = gdb.db.Updates(&book).Error
	if err != nil {
		return errors.New("unable to update book")
	}
	return nil

}
func (gdb *GormDB) DeleteBookOfUser(username string, id uint) error {
	_, err := gdb.GetBookById(username, id)
	if err != nil {
		return err
	}
	gdb.db.Delete(&Book{}, id)

	return nil
}
