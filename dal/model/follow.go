package model

type Follow struct {
	ID       int64 `gorm:"primarykey"`
	UserID   int64 `gorm:"index:idx_userid;not null" json:"user_id"`
	ToUserID int64 `gorm:"index:idx_userid;index:idx_userid_to;not null" json:"to_user_id"`
}
