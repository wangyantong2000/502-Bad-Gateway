package model

type Favorite struct {
	ID      int64 `gorm:"primarykey"`
	UserID  int64 `gorm:"index:idx_userid;not null" json:"user_id"`
	VideoID int64 `gorm:"index:idx_videoid;not null" json:"video_id"`
}
