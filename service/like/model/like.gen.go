// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameLike = "like"

// Like mapped from table <like>
type Like struct {
	ID       int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID   int64 `gorm:"column:user_id;not null" json:"user_id"`
	VideoID  int64 `gorm:"column:video_id;not null" json:"video_id"`
	Cancel   int32 `gorm:"column:cancel" json:"cancel"`
	AuthorID int64 `gorm:"column:author_id;not null" json:"author_id"`
}

// TableName Like's table name
func (*Like) TableName() string {
	return TableNameLike
}
