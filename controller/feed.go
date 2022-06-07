package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var videoList []Video

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed, pull the latest 30 videos uploaded by users
// need to consider to increase performance by concurrency
func Feed(c *gin.Context) {
	token, exist := c.GetQuery("token")
	Db.Transaction(func(tx *gorm.DB) error {

		err := tx.Table("videos").Order("id DESC").Limit(30).Find(&videoList).Error
		if err != nil {
			return err
		}
		if exist {
			var video *Video
			for i, _ := range videoList {
				video = &videoList[i]
				err = tx.Table("user_favorite_infos").
					Select("video_id").
					Where("token = ? AND video_id = ?", token, video.Id).
					Limit(1).Error
				if err != nil {
					return err
				}
				video.IsFavorite = true
			}
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videoList,
			NextTime:  time.Now().Unix(),
		})
		return nil
	})
}
