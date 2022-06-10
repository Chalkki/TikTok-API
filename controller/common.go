package controller

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64 `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Author        User  `gorm:"foreignKey:AuthorID;serializer:json" json:"author"`
	AuthorID      int64
	PlayUrl       string `json:"play_url" json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

type Comment struct {
	Id         int64  `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	User       User   `gorm:"serializer:json" json:"user"`
	UserId     int64  `json:"user_id"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
	VideoId    int64
	Video      Video `gorm:"foreignKey:VideoId" json:"-"`
}

type User struct {
	Id            int64  `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `gorm:"default:0" json:"follow_count"`
	FollowerCount int64  `gorm:"default:0" json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// UserLoginInfo Store the user login information
type UserLoginInfo struct {
	Token  string `json:"token"`
	UserId int64
	User   User `gorm:"foreignKey:UserId;serializer:json" json:"-"`
}

// UserFavoriteInfo Store the user favorite relation
type UserFavoriteInfo struct {
	UserId  int64
	VideoId int64
	Video   Video `gorm:"foreignKey:VideoId" json:"-"`
}

// UserFollowInfo Store the user follow information
type UserFollowInfo struct {
	UserId   int64
	ToUserId int64
}

//Scan and Value. It needs a scan method and a value method for gorm
//to serialize the customed data types.
func (user *User) Scan(value interface{}) error {
	val, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal([]byte(val), user)
}

func (user *User) Value() (driver.Value, error) {
	val, err := json.Marshal(&user)
	if err != nil {
		return nil, err
	}
	return val, nil
}
