package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// FavoriteAction
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	if token != loginInfo.Token {
		c.JSON(http.StatusOK,
			Response{StatusCode: 1,
				StatusMsg: "The provided token does not match to the user's login token"})
		return
	}
	// get the videoID from the request, and convert it from string to int64(if needed)
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK,
			Response{StatusCode: 1,
				StatusMsg: err.Error()})
		return
	}
	// get the action_type
	// action_type = 1 like
	// action_type = 2 unlike
	action_type := c.Query("action_type")
	// create the new userFavorite struct to store the information of favorite
	var userFavorite = UserFavoriteInfo{
		Token:   token,
		VideoId: videoID,
	}

	// get the reference of video struct in the video list
	// for add or subtract the favorite count
	var video *Video
	for _, v := range videoList {
		if videoID == v.Id {
			video = &v
			break
		}
	}

	//for more details, see the documentation of transaction in gorm
	//https://gorm.io/docs/transactions.html#content-inner
	Db.Transaction(func(tx *gorm.DB) error {
		// If action_type = "1", then we create a new record for the favorite relation.
		// If action_type != "1" (action_type == "2"), then we delete the record of the relation
		// between the specific user and video.
		// Set video.IsFavorite for current video list to make sure the state of favorite exists when
		// the user surf the video.
		if action_type == "1" {
			err = tx.Create(&userFavorite).Error
			if err != nil {
				c.JSON(http.StatusOK,
					Response{StatusCode: 1,
						StatusMsg: err.Error()})
				return err
			}
			video.FavoriteCount++
			video.IsFavorite = true

		} else {
			err = tx.
				Where("token = ? AND video_id = ?", userFavorite.Token, userFavorite.VideoId).
				Delete(&userFavorite).Error
			if err != nil {
				c.JSON(http.StatusOK,
					Response{StatusCode: 1,
						StatusMsg: err.Error()})
				return err
			}
			video.FavoriteCount--
			video.IsFavorite = false
		}

		tx.Model(&video).Update("favorite_count", video.FavoriteCount)
		if err != nil {
			c.JSON(http.StatusOK,
				Response{StatusCode: 1,
					StatusMsg: err.Error()})
			return err
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
		return nil
	})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	//c.JSON(http.StatusOK, VideoListResponse{
	//	Response: Response{
	//		StatusCode: 0,
	//	},
	//	VideoList: DemoVideos,
	//})
}
