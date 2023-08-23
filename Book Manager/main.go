package main

import (
	"Book_Manager/Authentication"
	"Book_Manager/Config"
	"Book_Manager/Database"
	"Book_Manager/Handlers"
	"Book_Manager/Validation"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

func main() {
	var cfg Config.Config
	logger := logrus.New()
	r := gin.Default()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetReportCaller(true)
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err.Error())
	}
	// Create a new instance of Database sql postgres
	gormDB, err := Database.CreateAndConnectToDb(cfg)
	if err != nil {
		logger.WithError(err).Panicln("can not connect to db")
	}
	if err := gormDB.CreateModel(); err != nil {
		logger.WithError(err).Fatalln("can not create table in db ")
	}
	logger.Infof("%+v", cfg)

	valid, err := Validation.CreateValidation(gormDB)
	if err != nil {
		logger.WithError(err).Fatal("can not Create instance od Validation ")

	}
	auth, err := Authentication.CreateAuthentication(gormDB, 10, logger)
	if err != nil {
		logger.WithError(err).Fatal("can not Create instance od Authentication ")
	}
	server := Handlers.BookMangerServer{
		Db:           gormDB,
		Logger:       logger,
		Authenticate: auth,
		Validation:   valid,
	}

	r.POST("/api/v1/auth/signup", server.HandleSignup)
	r.POST("/api/v1/auth/login", server.HandleLogin)
	r.POST("/api/v1/books", server.CreateBook)
	r.GET("/api/v1/books", server.GetAllBookOfUser)
	r.GET("/api/v1/books/:id", server.GetBookById)
	r.PUT("/api/v1/books/:id", server.UpdateBook)
	r.DELETE("/api/v1/books/:id", server.DeleteBook)

	if err := r.Run("localhost:8080"); err != nil {
		logrus.WithError(err).Fatalln("can not run server ")
	}

}
