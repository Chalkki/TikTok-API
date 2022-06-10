package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// FavoriteAction
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	var loginInfo UserLoginInfo
	err := loginInfo.GetUserInfo(token)
	if err != nil {
		c.JSON(http.StatusOK,
			Response{StatusCode: 1,
				StatusMsg: err.Error()})
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
		UserId:  loginInfo.UserId,
		VideoId: videoID,
	}

	// get the index of video struct in the video list
	// for add or subtract the favorite count
	var index int
	for i, _ := range feedVideoList {
		if videoID == feedVideoList[i].Id {
			index = i
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
			feedVideoList[index].FavoriteCount++
			feedVideoList[index].IsFavorite = true

		} else {
			err = tx.
				Where("user_id = ? AND video_id = ?", userFavorite.UserId, userFavorite.VideoId).
				Delete(&userFavorite).Error
			if err != nil {
				c.JSON(http.StatusOK,
					Response{StatusCode: 1,
						StatusMsg: err.Error()})
				return err
			}
			feedVideoList[index].FavoriteCount--
			feedVideoList[index].IsFavorite = false
		}
		var video Video
		tx.Model(&video).Where("id = ?", videoID).
			Update("favorite_count", feedVideoList[index].FavoriteCount)
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

// FavoriteList display the favorite list of the current user
func FavoriteList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: err.Error(),
		})
	}

	favoriteVideoList := make([]Video, 0)
	var userFavoriteInfos []UserFavoriteInfo

	Db.Transaction(func(tx *gorm.DB) error {
		err = tx.Table("user_favorite_infos").
			Where("user_id = ?", userId).
			Find(&userFavoriteInfos).
			Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else if err == nil {
			for i, _ := range userFavoriteInfos {
				var video Video
				tx.Where("id = ?", userFavoriteInfos[i].VideoId).First(&video)
				favoriteVideoList = append(favoriteVideoList, video)
			}
		}

		c.JSON(http.StatusOK, VideoListResponse{
			Response:  Response{StatusCode: 0},
			VideoList: favoriteVideoList,
		})
		return nil
	})
}
