package model

type Message struct {
	ID          int64  `gorm:"primarykey"`
	CreatedDate int64  `gorm:"index:idx_createdate;not null" json:"create_date"`
	UserID      int64  `gorm:"index:idx_userid;not null" json:"user_id"`
	ToUserID    int64  `gorm:"index:idx_userid;index:idx_userid_to;not null" json:"to_user_id"`
	Content     string `gorm:"type:varchar(255);not null" json:"content"`
}
