package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

//This is the loginInfo for the current user

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

//Register The function for a guest to create a new account
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := username + password

	//check whether there exists a user with the same username
	var existingUser User
	err := Db.Where("name = ?", username).First(&existingUser).Error
	//if not, it means that the new username is valid.
	//then the newUser information will be stored in the users table and the user_login_infos table.
	fmt.Println(err)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		newUser := User{
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
		}
		Db.Create(&newUser)

		newUserInfo := UserLoginInfo{
			User:   newUser,
			Token:  token,
			UserId: newUser.Id,
		}
		Db.Create(&newUserInfo)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   newUser.Id,
			Token:    username + password,
		})
	} else if err == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	}
}
func (i *UserLoginInfo) GetUserInfo(token string) error {
	err := Db.Where("token = ?", token).First(i).Error
	return err
}

// Login check if the account is stored in the user_login_infos table
// if yes, then we can allow the user to login
func Login(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	var loginInfo UserLoginInfo
	err := loginInfo.GetUserInfo(token)
	if err == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   loginInfo.User.Id,
			Token:    loginInfo.Token,
		})
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Something's wrong, please retry later"},
		})
	}

}

//UserInfo return the user information of the current login user
func UserInfo(c *gin.Context) {
	//check token with the tokens stored in the user_login_infos table.
	// if it exists in the table, return User to the server.
	// otherwise, return errors
	token := c.Query("token")
	userId := c.Query("user_id")
	var loginInfo UserLoginInfo
	var user User
	err := loginInfo.GetUserInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: err.Error(),
		})
		return
	}
	err = Db.Where("id = ?", userId).First(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     user,
	})
}
