package Handlers

import (
	"Book_Manager/Authentication"
	"Book_Manager/Database"
	"Book_Manager/Validation"
	"github.com/sirupsen/logrus"
)

type BookMangerServer struct {
	Db           *Database.GormDB
	Logger       *logrus.Logger
	Authenticate *Authentication.Authentication
	Validation   *Validation.Validation
}
