package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction Follow Action
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	var loginInfo UserLoginInfo
	if err := loginInfo.GetUserInfo(token); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	// this is the id of user that being followed
	toUserId, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	if loginInfo.UserId == toUserId {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Sorry, you cannot follow yourself."})
		return
	}
	//retrieve the information of current user and toUser
	var user User
	var toUser User
	tx := Db.Begin()
	if err = tx.Where("id = ?", loginInfo.UserId).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	if err = tx.Where("id = ?", toUserId).First(&toUser).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	record := UserFollowInfo{
		UserId:   loginInfo.UserId,
		ToUserId: toUserId,
	}
	//if action == "1", then add a new record in UserFollowInfo
	//the followerCount of the to_user would increment by one
	//the followCount of the current user would increment by one
	//
	//if action =="2", then delete the corresponding record in UserFollowInfo
	//the followerCount of the to_user would decrement by one
	//the followerCount of the current user would decrement by one
	action_type := c.Query("action_type")
	if action_type == "1" {
		if err = tx.Create(&record).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		user.FollowCount++
		toUser.FollowerCount++
	} else {
		if err = tx.Where("user_id = ? AND to_user_id = ?", record.UserId, record.ToUserId).
			Delete(&record).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		user.FollowCount--
		toUser.FollowerCount--
	}
	//update the user's followCount and toUser's followerCount
	if err = tx.Model(&user).Update("follow_count", user.FollowCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	if err = tx.Model(&toUser).Update("follower_count", toUser.FollowerCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	// update the author information to videos
	var video Video
	//serialize the toUser and user from json format to bytes
	toUserBytes, _ := toUser.Value()
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	userBytes, _ := user.Value()
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	if err = tx.Model(&video).Where("author_id = ?", toUserId).Update("author", toUserBytes).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	if err = tx.Model(&video).Where("author_id = ?", loginInfo.UserId).Update("author", userBytes).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// FollowList
// retrieve follow details from user_follow_infos
// return the users list according to follow details
func FollowList(c *gin.Context) {
	user_id := c.Query("user_id")
	var userFollowList []UserFollowInfo
	var userList = make([]User, 0)
	if err := Db.Where("user_id = ?", user_id).Find(&userFollowList).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	for i, _ := range userFollowList {
		var user User
		if err := Db.Where("id = ?", userFollowList[i].ToUserId).First(&user).Error; err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		userList = append(userList, user)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

// FollowerList similar to what FollowList does, but change the direction for query
func FollowerList(c *gin.Context) {
	user_id := c.Query("user_id")
	var userFollowerList []UserFollowInfo
	var userList = make([]User, 0)
	if err := Db.Where("to_user_id = ?", user_id).Find(&userFollowerList).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	for i, _ := range userFollowerList {
		var user User
		if err := Db.Where("id = ?", userFollowerList[i].UserId).First(&user).Error; err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		userList = append(userList, user)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}
