package controllers

import (
	"golang-base/forms"
	"golang-base/helpers"
	model "golang-base/models"
	"golang-base/services"

	"github.com/gin-gonic/gin"
)

var userModel = new(model.UserModel)

type UserController struct{}

func (u *UserController) Signup(c *gin.Context) {
	var data forms.SignupUsercommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide relevant fields"})
		c.Abort()
		return
	}

	result, _ := userModel.GetUserByEmail(data.Email)

	if result.Email != "" {
		c.JSON(403, gin.H{"message": "Email is already in use"})
		c.Abort()
		return
	}

	err := userModel.Signup(data)

	resetToken, _ := services.GenerateNonAuthToken(data.Email)

	link := "http://localhost:5000/api/v1/verify-account?verify_token=" + resetToken
	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"

	email := services.SendMail("Verify Account", body, data.Email, html, data.Name)

	if !email {
		c.JSON(500, gin.H{"message": "An issue occurent sending you an email"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem creating an account"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "New user account registered"})
}

func (u *UserController) Login(c *gin.Context) {
	var data forms.LoginUserCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required details"})
		c.Abort()
		return
	}

	result, err := userModel.GetUserByEmail(data.Email)

	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem logging into your account"})
		c.Abort()
		return
	}

	hashedPassword := []byte(result.Password)

	password := []byte(data.Password)

	err = helpers.PasswordCompare(password, hashedPassword)

	if err != nil {
		c.JSON(403, gin.H{"message": "Invalid user credentials"})
		c.Abort()
		return
	}

	jwtToken, refreshToken, err2 := services.GenerateToken(data.Email)

	if err2 != nil {
		c.JSON(403, gin.H{"message": "There was a problem logging you in, please try again"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Log in successful", "token": jwtToken, "refreshToken": refreshToken})
}

func (u *UserController) ResetLink(c *gin.Context) {
	var data forms.ResendCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required details"})
		c.Abort()
		return
	}

	result, err := userModel.GetUserByEmail(data.Email)

	if err != nil {
		c.JSON(500, gin.H{"message": "Something went wrong, try again later"})
		c.Abort()
		return
	}

	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	resetToken, _ := services.GenerateNonAuthToken(result.Email)

	link := "http://localhost:5000/api/v1/password-reset?reset_token=" + resetToken

	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"

	email := services.SendMail("Reset Password", body, result.Email, html, result.Name)

	// If email was sent, return 200 status code
	if email == true {
		c.JSON(200, gin.H{"messsage": "Check mail"})
		c.Abort()
		return
		// Return 500 status when something wrong happened
	} else {
		c.JSON(500, gin.H{"message": "An issue occured sending you an email"})
		c.Abort()
		return
	}
}

func (u *UserController) PasswordReset(c *gin.Context) {
	var data forms.PasswordResetCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required details"})
		c.Abort()
		return
	}

	if data.Password != data.Confirm {
		c.JSON(400, gin.H{"message": "Passwords do not match"})
		c.Abort()
		return
	}

	resetToken, _ := c.GetQuery("reset_token")

	userID, _ := services.DecodeNonAuthToken(resetToken)

	result, err := userModel.GetUserByEmail(userID)

	if err != nil {
		// Return response when we get an error while fetching user
		c.JSON(500, gin.H{"message": "Something wrong happened, try again later"})
		c.Abort()
		return
	}

	// Check if account exists
	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User accoun was not found"})
		c.Abort()
		return
	}

	newHashedPassword := helpers.GeneratePasswordHash([]byte(data.Password))

	_err := userModel.UpdateUserPass(userID, newHashedPassword)

	if _err != nil {
		// Return response if we are not able to update user password
		c.JSON(500, gin.H{"message": "Somehting happened while updating your password try again"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Password has been updated, log in"})
	c.Abort()
	return
}

func (u *UserController) VerifyLink(c *gin.Context) {
	var data forms.ResendCommand

	if (c.BindJSON(&data)) != nil {
		c.JSON(400, gin.H{"message": "Provided all fields"})
		c.Abort()
		return
	}

	result, err := userModel.GetUserByEmail(data.Email)

	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"message": "Something wrong happened, try again later"})
		c.Abort()
		return
	}

	resetToken, _ := services.GenerateNonAuthToken(result.Email)

	link := "http://localhost:5000/api/v1/verify-account?verify_token=" + resetToken
	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"

	email := services.SendMail("Verify Account", body, data.Email, html, result.Name)

	if email == true {
		c.JSON(200, gin.H{"message": "Check mail"})
		c.Abort()
		return
	} else {
		c.JSON(500, gin.H{"message": "An issue occurent sending you an email"})
		c.Abort()
		return
	}
}

func (u *UserController) VerifyAccount(c *gin.Context) {
	verifyToken, _ := c.GetQuery("verify_token")

	userID, _ := services.DecodeNonAuthToken(verifyToken)

	result, err := userModel.GetUserByEmail(userID)

	if err != nil {
		c.JSON(500, gin.H{"message": "Something wrong happened, try again later"})
		c.Abort()
		return
	}

	if result.Email == "" {
		c.JSON(404, gin.H{"mesasge": "User account was not found"})
		c.Abort()
		return
	}

	_err := userModel.VerifyAccount(userID)

	if _err != nil {
		c.JSON(500, gin.H{"message": "Something happened while verifying your account"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Account verified, log in"})
}

func (u *UserController) RefreshToken(c *gin.Context) {
	refreshToken := c.Request.Header["Refreshtoken"]

	if refreshToken == nil {
		c.JSON(403, gin.H{"message": "No refresh token provided"})
		c.Abort()
		return
	}

	email, err := services.DecodeRefreshToken(refreshToken[0])

	if err != nil {
		c.JSON(500, gin.H{"message": "Problem refreshing your session"})
		c.Abort()
		return
	}

	accessToken, _refreshToken, _err := services.GenerateToken(email)

	if _err != nil {
		c.JSON(500, gin.H{"message": "Problem creating new session"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Log in success", "token": accessToken, "refresh_token": _refreshToken})
}
