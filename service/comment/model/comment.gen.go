// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameComment = "comment"

// Comment mapped from table <comment>
type Comment struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID      int64     `gorm:"column:user_id;not null" json:"user_id"`
	VideoID     int64     `gorm:"column:video_id;not null" json:"video_id"`
	CommentText string    `gorm:"column:comment_text;not null" json:"comment_text"`
	CreatDate   time.Time `gorm:"column:creat_date;not null" json:"creat_date"`
	Cancel      int32     `gorm:"column:cancel" json:"cancel"`
}

// TableName Comment's table name
func (*Comment) TableName() string {
	return TableNameComment
}
