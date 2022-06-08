package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	var loginInfo UserLoginInfo
	err := loginInfo.GetUserInfo(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	var comment Comment
	var user User
	Db.Where("id = ?", loginInfo.UserId).First(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if actionType == "1" {

		content := c.Query("comment_text")
		now := time.Now().Unix()
		//Format YY-MM-DD
		unixT := time.Unix(now, 0)
		createDate := unixT.Format("06-01-02")
		comment = Comment{
			User:       user,
			Content:    content,
			CreateDate: createDate[3:],
			VideoId:    videoId,
		}
		//start of the transaction to create a new comment in comments table
		//the video's CommentCount would increment
		tx := Db.Begin()
		var video Video
		err = tx.First(&video).Where("id = ?", videoId).Error
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		err = tx.Model(&video).Update("comment_count", video.CommentCount+1).Where("id = ?", videoId).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		err = tx.Create(&comment).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		//end of the transaction
		c.JSON(http.StatusOK, CommentActionResponse{
			Response{StatusCode: 0},
			comment,
		})
	} else {
		commentId, _ := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		// start of the transaction to delete a comment in the comments table
		tx := Db.Begin()
		var video Video
		err = tx.First(&video).Where("id = ?", videoId).Error
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		err = tx.Model(&video).Update("comment_count", video.CommentCount+1).Where("id = ?", videoId).Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		err = tx.Where("id", commentId).Delete(&comment).Error
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		err = tx.Commit().Error
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
		//end of the transaction
		c.JSON(http.StatusOK, Response{
			StatusCode: 0, StatusMsg: "Delete successfully",
		})
	}

}

// CommentList each video will have their own comments' list in descending order
func CommentList(c *gin.Context) {
	var commentList []Comment
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	err := Db.Where("video_id = ?", videoId).Order("id DESC").Find(&commentList).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
}
