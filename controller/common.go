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
	Id            int64  `gorm:"primaryKey;autoIncrement:true" json:"id,omitempty"`
	Author        User   `gorm:"serializer:json" json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `gorm:"serializer:json" json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `gorm:"default:0" json:"follow_count"`
	FollowerCount int64  `gorm:"default:0" json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type UserLoginInfo struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
	User   User   `gorm:"foreignKey:UserId;serializer:json" json:"user"`
}

//we need a scan method and a value method for gorm
//to store the customed data types.
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
