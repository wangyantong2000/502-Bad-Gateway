package model

type Comment struct {
	ID          int64  `gorm:"primarykey"`
	CreatedTime int64  `gorm:"index;not null" json:"create_time"`
	UserID      int64  `gorm:"index:idx_userid;not null" json:"user_id"`
	VideoID     int64  `gorm:"index:idx_videoid;not null" json:"video_id"`
	Content     string `gorm:"type:varchar(255);not null" json:"content"`
}
