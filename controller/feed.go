package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

var feedVideoList []Video

// Feed pull the latest 30 videos uploaded on the oss server, and check whether the current liked the video
// Make sure to display properly the follow list and the favorite list when the feed is pulled by an existing user
// in the databse
func Feed(c *gin.Context) {
	token, exist := c.GetQuery("token")
	err := Db.Table("videos").Order("id DESC").Limit(30).Find(&feedVideoList).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: err.Error(),
		})
		return
	}

	Db.Transaction(func(tx *gorm.DB) error {
		//begin of the transaction to pull the feed
		err = tx.Table("videos").Order("id DESC").Limit(30).Find(&feedVideoList).Error
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return err
		}
		if exist {
			var userLoginInfo UserLoginInfo
			var userFollowInfos []UserFollowInfo
			var userFavoriteinfos []UserFavoriteInfo
			if err = tx.Where("token = ?", token).First(&userLoginInfo).Error; err != nil {
				return err
			}
			err = tx.Table("user_favorite_infos").
				Where("user_id = ?", userLoginInfo.UserId).
				Find(&userFavoriteinfos).
				Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
				return err
			}
			err = tx.Table("user_follow_infos").
				Where("user_id = ?", userLoginInfo.UserId).
				Find(&userFollowInfos).
				Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
				return err
			}
			for i, _ := range feedVideoList {
				var video *Video
				video = &feedVideoList[i]
				for j, _ := range userFavoriteinfos {
					if err == nil && userFavoriteinfos[j].VideoId == video.Id {
						video.IsFavorite = true
						break
					}
				}
				for j, _ := range userFollowInfos {
					if err == nil && userFollowInfos[j].ToUserId == video.AuthorID {
						video.Author.IsFollow = true
						break
					}
				}

			}
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: feedVideoList,
			NextTime:  time.Now().Unix(),
		})
		return nil
		//end of the transaction
	})
}
