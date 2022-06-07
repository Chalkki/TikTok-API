package controller

import (
	"errors"
	"fmt"
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
	fmt.Println(token)
	fmt.Println(exist)
	//err := Db.Table("videos").Order("id DESC").Limit(30).Find(&videoList).Error
	//if err != nil {
	//	c.JSON(http.StatusOK, FeedResponse{
	//		Response:  Response{StatusCode: 1, StatusMsg: err.Error()},
	//		VideoList: nil,
	//		NextTime:  time.Now().Unix(),
	//	})
	//	return
	//}
	//if exist {
	//	var video *Video
	//	var userFavoriteinfo UserFavoriteInfo
	//	for i, _ := range videoList {
	//		video = &videoList[i]
	//		err = Db.Table("user_favorite_infos").
	//			Where("token = ?", token).
	//			Limit(1).
	//			Find(&userFavoriteinfo).
	//			Error
	//		if err != nil {
	//			c.JSON(http.StatusOK, FeedResponse{
	//				Response:  Response{StatusCode: 1, StatusMsg: err.Error()},
	//				VideoList: nil,
	//				NextTime:  time.Now().Unix(),
	//			})
	//			return
	//		}
	//		video.IsFavorite = true
	//	}
	//}

	Db.Transaction(func(tx *gorm.DB) error {

		err := tx.Table("videos").Order("id DESC").Limit(30).Find(&videoList).Error
		if err != nil {
			return err
		}
		if exist {
			var video *Video
			var userFavoriteinfo UserFavoriteInfo
			for i, _ := range videoList {
				video = &videoList[i]
				err = tx.Table("user_favorite_infos").
					Where("token = ?", token).
					Where("video_id = ?", video.Id).
					Limit(1).
					Find(&userFavoriteinfo).
					Error
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					fmt.Println(err.Error())
					return err
				}
				// in case the table is empty, we need check whether the query returns empty struct
				// if the record could be found in the userFavoriteInfos, we can see the video.IsFavorite to true
				if err == nil && userFavoriteinfo.VideoId != 0 {
					video.IsFavorite = true
				}
			}
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videoList,
			NextTime:  time.Now().Unix(),
		})
		return nil
	})

	//c.JSON(http.StatusOK, FeedResponse{
	//	Response:  Response{StatusCode: 0},
	//	VideoList: videoList,
	//	NextTime:  time.Now().Unix(),
	//})
}
