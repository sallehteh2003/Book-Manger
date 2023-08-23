package Handlers

import (
	"Book_Manager/Authentication"
	"Book_Manager/Database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type signupRequestBody struct {
	Username    string `json:"user_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Gender      string `json:"gender" binding:"required"`
}

type loginRequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *BookMangerServer) HandleSignup(c *gin.Context) {
	var reqData signupRequestBody

	// unmarshal json
	err := c.BindJSON(&reqData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	//Validate user data
	if err := s.Validation.ValidateData(reqData.Email, reqData.PhoneNumber); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid email or phone number"})
		return
	}

	// check user duplicate
	if err := s.Db.CheckUserDuplicate(reqData.Username, reqData.PhoneNumber, reqData.Email); err == nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "user already exist"})
		return
	}

	// Create new user in database
	user := &Database.UserDb{

		Username:    reqData.Username,
		Firstname:   reqData.FirstName,
		Lastname:    reqData.LastName,
		Email:       reqData.Email,
		Password:    reqData.Password,
		PhoneNumber: reqData.PhoneNumber,
		Gender:      reqData.Gender,
	}
	if err := s.Db.CreateNewUser(user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "can not Create new user "})
		return
	}
	c.IndentedJSON(http.StatusCreated, "")
	return
}
func (s *BookMangerServer) HandleLogin(c *gin.Context) {
	var reqData loginRequestBody

	// unmarshal json
	err := c.BindJSON(&reqData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	// Authenticate of user
	err = s.Authenticate.AuthenticateUserWithCredentials(Authentication.Credentials{
		Username: reqData.Username,
		Password: reqData.Password,
	})
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	// Create the JWT token
	token, err := s.Authenticate.GenerateJwtToken(reqData.Username)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "can not Create token"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"access_token": *token})
	return

}
