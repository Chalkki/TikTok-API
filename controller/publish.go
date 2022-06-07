package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	fmt.Println("the provided token is ", token)
	fmt.Println("the login token is", loginInfo.Token)
	if token != loginInfo.Token {
		c.JSON(http.StatusOK, Response{StatusCode: 1,
			StatusMsg: "The provided token does not match with the current token"})
		return
	}

	// read the raw video data from the http request
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// save the video in the ./public path
	filename := filepath.Base(data.Filename)
	user := loginInfo.User
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// See the final id of the video table, and add 1 to the final id as the id of new video
	var lastVideo Video
	err = Db.Select("id").Last(&lastVideo).Error
	var videoID int64
	if errors.Is(err, gorm.ErrRecordNotFound) {
		videoID = 1
		err = nil
	} else if err != nil {
		return
	}
	videoID = lastVideo.Id + 1
	err = nil
	// upload the video to the oos cloud first
	videoName := strconv.FormatInt(videoID, 10) + "-video.mp4"
	snapShotName := strconv.FormatInt(videoID, 10) + "-cover.jpeg"
	if err = UploadVideo(videoName, snapShotName, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title := c.PostForm("title")
	// The ObjectURLPrefix is the prefix of the PlayUrl and CoverUrl addresses
	// e.g. "https://www.w3schools.com/html/" is the URLPrefix of the demo video's playURL
	//update newVideo to the videos' table in the database
	var newVideo = Video{
		Author:   loginInfo.User,
		PlayUrl:  ObjectURLPrefix + videoName,
		CoverUrl: ObjectURLPrefix + snapShotName,
		Title:    title,
		AuthorID: loginInfo.UserId,
	}
	err = Db.Create(&newVideo).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	var videoList []Video
	userID := c.Query("user_id")
	token := c.Query("token")

	// search for the token in the loginInfo table that matched to the token provided by the client.
	// if the token doesn't match to each other, the server will inform an error to client.
	err := Db.Where("token = ?", token).First(&loginInfo).Error
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User not verified, please login again and try",
			},
		})
		return
	}

	// Search for all the videos uploaded by the user
	//in the videos' table in descending way and provide them to the client
	err = Db.Table("videos").
		Where("author_id = ?", userID).
		Order("id DESC").
		Find(&videoList).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
